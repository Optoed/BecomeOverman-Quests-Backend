-- Удаление таблиц, если они существуют (с правильным порядком и CASCADE)
DROP TABLE IF EXISTS user_achievements CASCADE;
DROP TABLE IF EXISTS achievements CASCADE;
DROP TABLE IF EXISTS user_coin_transactions CASCADE;
DROP TABLE IF EXISTS user_daily_streaks CASCADE;
DROP TABLE IF EXISTS user_tasks CASCADE;
DROP TABLE IF EXISTS task_variants CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS user_completed_quests CASCADE;
DROP TABLE IF EXISTS user_current_quests CASCADE;
DROP TABLE IF EXISTS user_quests CASCADE;
DROP TABLE IF EXISTS quest_tasks CASCADE;
DROP TABLE IF EXISTS quests CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS friends CASCADE;
DROP TABLE IF EXISTS shared_quests CASCADE;

-- Удаление типов
DROP TYPE IF EXISTS category_name CASCADE;
DROP TYPE IF EXISTS difficulty_level CASCADE;
DROP TYPE IF EXISTS task_type CASCADE;
DROP TYPE IF EXISTS rarity CASCADE;

-- Создание ENUM типов
CREATE TYPE category_name AS ENUM ('health', 'intelligence', 'charisma', 'willpower');
CREATE TYPE rarity AS ENUM ('free', 'common', 'rare', 'epic', 'legendary');
CREATE TYPE task_type AS ENUM ('daily', 'weekly', 'special', 'user_generated');

-- Таблица пользователей с расширенными полями
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    xp_points INT DEFAULT 0,
    coin_balance INT DEFAULT 0,
    level INT DEFAULT 1,

    -- TODO: если появятся еще ветки то поправим
    health_level INT DEFAULT 0,
    mental_health_level INT DEFAULT 0,
    intelligence_level INT DEFAULT 0,
    charisma_level INT DEFAULT 0,
    willpower_level INT DEFAULT 0,

    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица задач
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,

    difficulty INT DEFAULT 0,
    rarity VARCHAR(255) NOT NULL DEFAULT 'free', -- 'free', 'common', 'rare', 'epic', 'legendary'
    category VARCHAR(255) NOT NULL,
    base_xp_reward INT NOT NULL DEFAULT 0,
    base_coin_reward INT NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    -- type task_type NOT NULL DEFAULT 'daily',
    -- cooldown_hours INT DEFAULT 24,
);

-- Транзакции валюты пользователя
CREATE TABLE user_coin_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    reference_id INT, -- ID связанной сущности (задача, покупка и т.д.)

    transaction_type VARCHAR(50) NOT NULL, -- 'earned', 'spent', 'bonus'
    amount INT NOT NULL,
    
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Достижения
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    criteria_json JSONB NOT NULL, -- например: {"tasks_completed": 100}
    bonus_json JSONB, -- например заморозка или пассивный бонус
    reward_xp INT DEFAULT 0,
    reward_coin INT DEFAULT 0,
    is_secret BOOLEAN DEFAULT FALSE
);

-- Достижения пользователей
CREATE TABLE user_achievements (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    achievement_id INT NOT NULL,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_achievement UNIQUE (user_id, achievement_id)
);

-- Таблица квестов
CREATE TABLE quests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255) NOT NULL, -- 'health', 'willpower', 'intelligence', 'charisma'
    rarity VARCHAR(255) NOT NULL, -- 'free', 'common', 'rare', 'epic', 'legendary'
    difficulty INT NOT NULL DEFAULT 0,
    price INT NOT NULL DEFAULT 0,
    tasks_count INT DEFAULT 1, -- Сколько задач в квесте
    conditions_json JSONB,
    bonus_json JSONB,
    is_sequential BOOLEAN DEFAULT FALSE,   -- Нужно ли выполнять по порядку
    reward_xp INT NOT NULL,
    reward_coin INT NOT NULL,
    time_limit_hours INT DEFAULT 0          -- Ограничение по времени (опционально)
);


-- задачи пользователя
CREATE TABLE user_tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    quest_id INT REFERENCES quests(id) ON DELETE SET NULL,  -- опционально

    status VARCHAR(50) NOT NULL DEFAULT 'active',         -- not_started, active, completed, TODO: failed
    scheduled_start TIMESTAMP,
    scheduled_end TIMESTAMP,
    deadline TIMESTAMP,
    duration INT,                                        -- время выделенное на задачу в минутах (?)
    updated_by_ai BOOLEAN DEFAULT FALSE,                 -- была ли запланирована AI

    is_confirmed BOOL DEFAULT FALSE NOT NULL,            -- прежнее поле
    completed_at TIMESTAMP,                              -- прежнее поле
    xp_gained INT NOT NULL DEFAULT 0,
    coin_gained INT NOT NULL DEFAULT 0,

    CONSTRAINT unique_user_task UNIQUE (user_id, task_id)
);

-- Связь квестов и задач (какие задачи входят в квест)
CREATE TABLE quest_tasks (
    id SERIAL PRIMARY KEY,
    quest_id INT NOT NULL,
    task_id INT NOT NULL,

    task_order INT,                        -- Порядок (если is_sequential = TRUE)

    FOREIGN KEY (quest_id) REFERENCES quests(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- Прогресс пользователя по квестам
CREATE TABLE user_quests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    quest_id INT NOT NULL,

    status VARCHAR(255) NOT NULL DEFAULT 'purchased', -- "purchased", "started", "failed", "completed"

    xp_gained INT,
    coin_gained INT,

    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    expires_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (quest_id) REFERENCES quests(id) ON DELETE CASCADE
);

-- friends
-- Друзья
CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, friend_id)
);

-- Совместные квесты
CREATE TABLE shared_quests (
    id SERIAL PRIMARY KEY,
    quest_id INTEGER NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    user1_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);