## 0) Подготовка

Открыть PowerShell в папке проекта:

```powershell
cd C:\Users\user\BecomeOverman-Quests-Backend-FORK
```

Запустить сервисы:

```powershell
docker compose up -d --build
docker compose --profile consumer up -d --build consumer
```

Проверить статус:

```powershell
docker compose ps
```

Ожидаем `Up` для `backend`, `kafka`, `kafka-ui`, `postgres`, `consumer`.

---

## 1) Что показать в браузере

- `http://localhost:8080/` — backend API жив.
- `http://localhost:8081/` — Kafka UI.

В Kafka UI показать:
- `Online clusters: 1`
- топики:
  - `becomeoverman.user.events`
  - `becomeoverman.quest.events`

---

## 2) Логин и получение JWT

```powershell
$loginBody = @{
  username = "kafka_user_1"
  password = "12345678"
} | ConvertTo-Json

$loginResp = Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/auth/login" `
  -ContentType "application/json" `
  -Body $loginBody

$token = $loginResp.token
$headers = @{ Authorization = "Bearer $token" }

$loginResp
```

---

## 3) Сценарий A: покупка квеста -> user.quest_purchased

### 3.1 Получить квест из магазина

```powershell
$shop = Invoke-RestMethod -Method Get `
  -Uri "http://localhost:8080/quests/shop" `
  -Headers $headers

$shop
$questId = $shop.id
```

### 3.2 Купить квест

```powershell
$purchaseBody = @{ status = "purchased" } | ConvertTo-Json

Invoke-RestMethod -Method Patch `
  -Uri "http://localhost:8080/users/me/quests/$questId" `
  -Headers $headers `
  -ContentType "application/json" `
  -Body $purchaseBody
```

### 3.3 Показать событие в consumer

```powershell
docker compose logs --tail=120 consumer
```

Искать в выводе:
- `topic=becomeoverman.user.events`
- `event_type":"user.quest_purchased"`

---

## 4) Сценарий B: создание квеста -> quest.created

### 4.1 Создать квест через dev endpoint

```powershell
$devCreateBody = @{
  quest = @{
    title = "Kafka Created Quest"
    description = "Created via dev endpoint for kafka test"
    category = "health"
    rarity = "free"
    difficulty = 1
    price = 0
    tasks_count = 1
    reward_xp = 10
    reward_coin = 5
    time_limit_hours = 24
    is_sequential = $false
  }
  tasks = @(
    @{
      title = "Kafka Created Task"
      description = "Task for quest.created event"
      difficulty = 1
      rarity = "free"
      category = "health"
      base_xp_reward = 3
      base_coin_reward = 1
      task_order = 1
    }
  )
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/quests/dev-create" `
  -Headers $headers `
  -ContentType "application/json" `
  -Body $devCreateBody
```

### 4.2 Показать событие в consumer

```powershell
docker compose logs --tail=150 consumer
```

Искать в выводе:
- `topic=becomeoverman.quest.events`
- `event_type":"quest.created"`

---

## 5) JWT истек / Missing or invalid token (быстрый фикс)

Если любая защищенная команда возвращает:

`{"error":"Missing or invalid token"}`

просто выполнить этот блок заново:

```powershell
$loginBody = @{
  username = "kafka_user_1"
  password = "12345678"
} | ConvertTo-Json

$loginResp = Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/auth/login" `
  -ContentType "application/json" `
  -Body $loginBody

$token = $loginResp.token
$headers = @{ Authorization = "Bearer $token" }
```

После этого повторить исходную команду.

---

## 6) Остановка после демо

Остановить всё:

```powershell
docker compose down
```

Полный сброс с данными БД:

```powershell
docker compose down -v
```
