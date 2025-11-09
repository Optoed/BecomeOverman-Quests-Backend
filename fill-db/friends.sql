-- Создаем первый квест с задачами
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Утренний дружеский марафон', 
        'Совместный недельный челлендж для развития силы воли и здоровья', 
        'health', 
        'common', 
        2, 
        15, 
        5,
        150, 
        75, 
        168
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Пробуждение в 6 утра', 'Проснуться в 6 утра без будильника', 1, 'common', 'willpower', 20, 10),
    ('Утренняя зарядка', 'Сделать 15-минутную зарядку', 1, 'common', 'health', 25, 12),
    ('Совместная пробежка', 'Пробежать 3 км вместе с другом', 2, 'common', 'health', 30, 15),
    ('Медитация 10 минут', 'Помедитировать 10 минут утром', 1, 'common', 'mental_health', 20, 10),
    ('План на день', 'Составить план задач на день', 1, 'common', 'intelligence', 15, 8)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;

-------------------------------------------------------------------------------------------------------------------
-- Квест 3: Базовое саморазвитие (бесплатный)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Основы продуктивности', 
        'Недельный план для формирования полезных привычек', 
        'willpower', 
        'common', 
        1, 
        0,  -- бесплатный
        4,
        80, 
        40, 
        168
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Утренний ритуал', 'Выполнить 3 запланированных утренних действия', 1, 'common', 'willpower', 15, 8),
    ('Фокус-сессия', 'Работать без отвлечений 45 минут', 2, 'common', 'intelligence', 20, 10),
    ('Физическая активность', 'Сделать 15 минут любой физической активности', 1, 'common', 'health', 15, 7),
    ('Вечерний анализ', 'Записать 3 достижения за день', 1, 'common', 'mental_health', 10, 5)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;

-------------------------------------------------------------------------------------------------------------------

-- Квест 4: Дуэт креативности (стоимость 100)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Творческий дуэт', 
        'Совместное создание творческого проекта за неделю', 
        'charisma', 
        'rare', 
        2, 
        100, 
        5,
        180, 
        90, 
        168
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Мозговой штурм', 'Провести совместный мозговой штурм идей', 1, 'common', 'intelligence', 25, 12),
    ('Разработка концепции', 'Создать детальный план проекта', 2, 'common', 'intelligence', 30, 15),
    ('Распределение ролей', 'Определить зоны ответственности каждого', 1, 'common', 'charisma', 20, 10),
    ('Совместная работа', 'Потратить 2 часа на совместное создание', 3, 'rare', 'willpower', 40, 20),
    ('Презентация проекта', 'Показать результат друг другу', 2, 'common', 'charisma', 25, 13)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;

-------------------------------------------------------------------------------------------------------------------

-- Квест 6: Духовное развитие (бесплатный)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Путь к гармонии', 
        '14-дневный путь к внутреннему балансу и осознанности', 
        'mental_health', 
        'common', 
        2, 
        0,  -- бесплатный
        5,
        120, 
        60, 
        336  -- 14 дней
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Утренняя медитация', 'Помедитировать 10 минут после пробуждения', 1, 'common', 'mental_health', 15, 8),
    ('Дневник благодарности', 'Записать 5 вещей, за которые вы благодарны', 1, 'common', 'mental_health', 12, 6),
    ('Цифровой детокс', 'Провести 4 часа без гаджетов', 2, 'common', 'willpower', 25, 12),
    ('Прогулка на природе', 'Погулять в парке или лесу 1 час', 1, 'common', 'health', 18, 9),
    ('Анализ прогресса', 'Проанализировать изменения за 2 недели', 2, 'common', 'intelligence', 20, 10)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;

-------------------------------------------------------------------------------------------------------------------

-- Квест 7: Кулинарный вызов с партнером (стоимость 1)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Кулинарный дуэт', 
        'Неделя совместного кулинарного мастерства', 
        'charisma', 
        'rare', 
        3, 
        1,  -- символическая стоимость
        4,
        160, 
        80, 
        168
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Планирование меню', 'Составить совместное меню на неделю', 2, 'common', 'intelligence', 30, 15),
    ('Совместная готовка', 'Приготовить сложное блюдо вместе', 3, 'rare', 'charisma', 45, 22),
    ('Обмен рецептами', 'Научить партнера своему фирменному рецепту', 2, 'common', 'charisma', 35, 18),
    ('Гастрономический вечер', 'Устроить ужин с приготовленными блюдами', 2, 'common', 'charisma', 30, 15)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;

-------------------------------------------------------------------------------------------------------------------

-- Квест 8: Фитнес-марафон (бесплатный)
WITH new_quest AS (
    INSERT INTO quests (
        title, description, category, rarity, difficulty, price, tasks_count,
        reward_xp, reward_coin, time_limit_hours
    ) VALUES (
        'Фитнес-марафон', 
        '21 день для формирования спортивной привычки', 
        'health', 
        'common', 
        3, 
        0,  -- бесплатный
        7,
        210, 
        105, 
        504  -- 21 день
    ) RETURNING id
),
new_tasks AS (
    INSERT INTO tasks (title, description, difficulty, rarity, category, base_xp_reward, base_coin_reward) VALUES 
    ('Первая тренировка', 'Выполнить базовый комплекс упражнений', 1, 'common', 'health', 20, 10),
    ('Кардио-день', '30 минут кардио-нагрузки', 2, 'common', 'health', 25, 12),
    ('Силовая тренировка', 'Проработать основные группы мышц', 3, 'common', 'health', 30, 15),
    ('День восстановления', 'Растяжка и легкая активность', 1, 'common', 'health', 15, 8),
    ('Интервальная тренировка', 'Высокоинтенсивный интервальный тренинг', 4, 'rare', 'health', 35, 18),
    ('Полный комплекс', 'Выполнить полную программу тренировки', 3, 'common', 'willpower', 30, 15),
    ('Финальное испытание', 'Показать лучший результат в тестовом упражнении', 4, 'rare', 'health', 35, 17)
    RETURNING id
),
numbered_tasks AS (
    SELECT id, ROW_NUMBER() OVER () as task_order
    FROM new_tasks
)
INSERT INTO quest_tasks (quest_id, task_id, task_order)
SELECT (SELECT id FROM new_quest), id, task_order
FROM numbered_tasks;