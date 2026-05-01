# Become Overman — Quests Backend

<div align="center">

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge\&logo=go\&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-HTTP%20API-00ADD8?style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=for-the-badge\&logo=postgresql\&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge\&logo=jsonwebtokens\&logoColor=white)
![AI](https://img.shields.io/badge/AI-Quest%20Generation-7C3AED?style=for-the-badge)

**Backend-сервис для Become Overman: геймифицированной системы саморазвития с квестами, прокачкой характеристик, внутриигровой экономикой, AI-планированием и интеграцией с рекомендательным микросервисом.**

[Recommendation Service](https://github.com/Optoed/BecomeOverman-RecommendationSystem) · [API](#api) · [Architecture](#архитектура) · [Database](#база-данных)

</div>

---

## О проекте

**Become Overman** — это дипломный проект: геймифицированная платформа саморазвития, где личные цели превращаются в квесты, задачи, уровни, награды и прогресс по направлениям развития.

Этот репозиторий содержит основной backend на Go. Он отвечает за пользователей, авторизацию, магазин квестов, покупку и выполнение квестов, прогресс пользователя, внутриигровую валюту, совместные квесты, AI-генерацию и интеграцию с отдельным сервисом рекомендаций.

Проект выполнен в рамках обучения в **СГУ КНиИТ**.

---

## Идея

Большинство habit tracker и todo-приложений быстро превращаются в скучный список галочек. Become Overman строится вокруг другой логики:

> пользователь не просто закрывает задачи, а проходит квесты, получает XP, монеты, прокачивает характеристики и открывает новые испытания.

Квест может иметь цену, редкость, сложность, набор задач, ограничение по времени, награды и условия доступности. Это делает саморазвитие ближе к RPG-прогрессии: есть магазин, инвентарь, активные испытания, история прохождения и персональные рекомендации.

---

## Что реализовано

* регистрация и авторизация пользователей;
* JWT-based authentication;
* хранение пользователей, квестов, задач и прогресса в PostgreSQL;
* магазин доступных квестов;
* покупка квеста за внутриигровую валюту;
* запуск купленного квеста;
* выполнение задач внутри квеста;
* завершение квеста и начисление наград;
* учет XP, монет, уровней и характеристик пользователя;
* категории развития: `health`, `intelligence`, `charisma`, `willpower`;
* редкости квестов: `free`, `common`, `rare`, `epic`, `legendary`;
* shared quests для совместного прохождения;
* friends-модель для социальных механик;
* AI-генерация квестов;
* AI-планирование расписания;
* поиск квестов через recommendation service;
* рекомендации квестов;
* рекомендации пользователей / потенциальных друзей.

---

## Архитектура

```mermaid
flowchart TD
    A["Client / Frontend"] --> B["Go Backend API"]

    B --> C["Handlers"]
    C --> D["Services"]
    D --> E["Repositories"]
    E --> F["PostgreSQL"]

    D --> G["AI integrations"]
    D --> H["Recommendation Service"]

    H --> I["Semantic quest search"]
    H --> J["Quest recommendations"]
    H --> K["User recommendations"]

    B --> L["JWT middleware"]
    B --> M["CORS"]
```

Backend построен в классической многослойной архитектуре:

```text
HTTP handlers → services → repositories → PostgreSQL
```

Такой подход отделяет API-слой от бизнес-логики и SQL-доступа, а интеграции с AI и рекомендательным сервисом остаются внутри service layer.

---

## Стек

| Область             | Технологии                              |
| ------------------- | --------------------------------------- |
| Язык                | Go                                      |
| HTTP framework      | Gin                                     |
| Database            | PostgreSQL                              |
| DB access           | sqlx, lib/pq                            |
| Auth                | JWT                                     |
| Config              | godotenv, env variables                 |
| Migrations / schema | SQL, golang-migrate dependency          |
| Architecture        | handlers / services / repositories      |
| External services   | Recommendation service, AI integrations |

---

## База данных

Основные сущности:

| Таблица                  | Назначение                                                                |
| ------------------------ | ------------------------------------------------------------------------- |
| `users`                  | пользователи, XP, монеты, уровень, характеристики, streak                 |
| `quests`                 | квесты: категория, редкость, сложность, цена, награды, условия            |
| `tasks`                  | отдельные задания внутри квестов                                          |
| `quest_tasks`            | связь квестов и задач, порядок выполнения                                 |
| `user_quests`            | состояние квеста у пользователя: purchased / started / failed / completed |
| `user_tasks`             | задачи пользователя, расписание, deadline, AI-планирование                |
| `user_coin_transactions` | история начисления и списания монет                                       |
| `achievements`           | достижения и бонусы                                                       |
| `user_achievements`      | полученные достижения пользователя                                        |
| `friends`                | социальные связи пользователей                                            |
| `shared_quests`          | совместные квесты                                                         |

---

## Quest lifecycle

```mermaid
stateDiagram-v2
    [*] --> Available
    Available --> Purchased: buy quest
    Purchased --> Started: start quest
    Started --> TaskProgress: complete tasks
    TaskProgress --> Completed: complete quest
    TaskProgress --> Failed: time limit / fail condition
    Failed --> Purchased: retry
    Completed --> [*]
```

Типовой сценарий:

1. Пользователь получает список доступных квестов.
2. Backend проверяет условия: уровень, монеты, уже купленные/завершенные квесты.
3. Пользователь покупает квест.
4. Квест попадает в список пользовательских квестов, но таймер еще не стартует.
5. Пользователь запускает квест.
6. Выполняет задачи внутри квеста.
7. Завершает квест и получает XP, монеты и прогресс.
8. Данные пользователя отправляются в recommendation service для обновления профиля рекомендаций.

---

## API

### Auth / Users

| Method | Endpoint         | Назначение                    |
| ------ | ---------------- | ----------------------------- |
| `POST` | `/auth/register` | регистрация пользователя      |
| `POST` | `/auth/login`    | авторизация и получение JWT   |
| `GET`  | `/users/me`      | профиль текущего пользователя |

### Quests

| Method | Endpoint               | Назначение                                               |
| ------ | ---------------------- | -------------------------------------------------------- |
| `GET`  | `/quests/available`    | доступные квесты                                         |
| `GET`  | `/quests/shop`         | магазин квестов                                          |
| `GET`  | `/quests/search?q=...` | семантический поиск квестов через recommendation service |
| `GET`  | `/quests/:questID`     | детали квеста                                            |
| `POST` | `/quests`              | AI-генерация квеста                                      |
| `POST` | `/quests/shared`       | создание совместного квеста                              |

### User quests

| Method  | Endpoint                                  | Назначение                        |
| ------- | ----------------------------------------- | --------------------------------- |
| `GET`   | `/users/me/quests`                        | все квесты пользователя           |
| `GET`   | `/users/me/quests?status=active`          | активные квесты                   |
| `GET`   | `/users/me/quests?status=completed`       | завершенные квесты                |
| `PATCH` | `/users/me/quests/:questID`               | купить / начать / завершить квест |
| `PATCH` | `/users/me/quests/:questID/tasks/:taskID` | отметить задачу выполненной       |

### Recommendations

| Method | Endpoint                            | Назначение                        |
| ------ | ----------------------------------- | --------------------------------- |
| `GET`  | `/users/me/recommendations/quests`  | персональные рекомендации квестов |
| `GET`  | `/users/me/recommendations/friends` | рекомендации пользователей        |

### AI schedule

| Method | Endpoint     | Назначение                    |
| ------ | ------------ | ----------------------------- |
| `POST` | `/schedules` | генерация расписания через AI |

---

## Интеграция с Recommendation Service

Backend связан с отдельным Python/FastAPI-сервисом:

[BecomeOverman-RecommendationSystem](https://github.com/Optoed/BecomeOverman-RecommendationSystem)

Он отвечает за:

* семантический поиск квестов;
* поиск похожих квестов;
* персональные рекомендации;
* рекомендации пользователей по похожести профилей;
* хранение эмбеддингов квестов и пользователей.

После покупки квеста backend асинхронно обновляет профиль пользователя в recommendation service, чтобы рекомендации учитывали новые действия пользователя.

---

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/Optoed/BecomeOverman-Quests-Backend.git
cd BecomeOverman-Quests-Backend
```

### 2. Настроить `.env`

```bash
cp env.example .env
```

Пример:

```env
DATABASE_URL=postgresql://postgres:mysecretpassword@localhost:5432/mydb?sslmode=disable
JWT_SECRET=your_super_secret_key
TOKEN_EXPIRE_HOURS=24
```

### 3. Подготовить PostgreSQL

Создать базу данных и применить схему:

```bash
psql -U postgres -d mydb -f initDB.sql
```

При необходимости загрузить тестовые данные:

```bash
psql -U postgres -d mydb -f fillDB.sql
```

### 4. Запустить backend

```bash
go run ./cmd/app
```

По умолчанию API запускается на:

```text
http://localhost:8080
```

---

## Связанные репозитории

* [BecomeOverman-RecommendationSystem](https://github.com/Optoed/BecomeOverman-RecommendationSystem) — Python/FastAPI микросервис семантического поиска и рекомендаций.

---

## Что можно улучшить дальше

* добавить Swagger/OpenAPI-документацию;
* вынести все интеграционные URL в env/config;
* добавить Docker Compose для backend + PostgreSQL + recommendation service;
* покрыть бизнес-логику unit/integration-тестами;
* формализовать статусы квестов через enum/const;
* добавить миграции вместо ручного SQL-запуска;
* усилить валидацию AI-generated квестов;
* добавить observability: structured logs, metrics, tracing.

---

## Авторство

Проект выполнен в рамках дипломной работы студентом 4-ого курса **СГУ КНиИТ** Стеклянниковым Петром Сергеевичем.

---

## Лицензия

MIT License.
