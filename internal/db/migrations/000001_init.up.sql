CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    passport_number VARCHAR(11) NOT NULL UNIQUE,
    surname VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50),
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') NOT NULL
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE tasks (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid UUID UNIQUE REFERENCES users(uuid) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') NOT NULL,
    end_time TIMESTAMPTZ
);

CREATE TABLE task_histories (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL
);

-- Insert 10 users
INSERT INTO users (passport_number, surname, name, patronymic, address)
VALUES
('8234 557891', 'Ivanov', 'Ivan', 'Ivanovich', 'ул. Ленина, д. 1, кв. 1, Москва'),
('1434 564890', 'Petrov', 'Petr', 'Petrovich', 'пр. Мира, д. 2, кв. 10, Санкт-Петербург'),
('1534 567890', 'Sidorov', 'Sidr', NULL, 'ул. Пушкина, д. 3, кв. 20, Новосибирск'),
('1239 565890', 'Smirnov', 'Sergey', 'Sergeevich', 'ул. Чайковского, д. 4, кв. 30, Екатеринбург'),
('1274 567890', 'Kuznetsov', 'Nikolay', NULL, 'ул. Советская, д. 5, кв. 40, Казань'),
('6234 167890', 'Popov', 'Alexey', 'Alexeevich', 'ул. Лермонтова, д. 6, кв. 50, Нижний Новгород'),
('1234 567896', 'Vasiliev', 'Vasiliy', NULL, 'ул. Горького, д. 7, кв. 60, Самара'),
('4234 561898', 'Pavlov', 'Pavel', 'Pavlovich', 'ул. Маяковского, д. 8, кв. 70, Омск'),
('2234 567197', 'Rybakov', 'Roman', NULL, 'ул. Есенина, д. 9, кв. 80, Челябинск'),
('5234 567823', 'Kovalev', 'Konstantin', 'Konstantinovich', 'ул. Некрасова, д. 10, кв. 90, Ростов-на-Дону');

-- Add current tasks for each user
INSERT INTO tasks (user_uuid, name, start_time)
SELECT
    uuid,
    'Current Task ' || row_number() OVER (ORDER BY uuid),
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + (interval '1 minute' * floor(random() * 60))
FROM users;

-- Add completed tasks for each user
WITH task_intervals AS (
    SELECT interval '1 day' * 1 AS task_interval UNION ALL
    SELECT interval '1 day' * 2 UNION ALL
    SELECT interval '1 week' * 1 UNION ALL
    SELECT interval '1 week' * 2 UNION ALL
    SELECT interval '1 month' * 1 UNION ALL
    SELECT interval '1 month' * 2 UNION ALL
    SELECT interval '1 year' * 1 UNION ALL
    SELECT interval '1 year' * 2
)
INSERT INTO task_histories (user_uuid, name, start_time, end_time)
SELECT
    u.uuid,
    'Completed Task ' || row_number() OVER (PARTITION BY u.uuid ORDER BY t.task_interval),
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC' - t.task_interval + (interval '1 minute' * floor(random() * 60)),
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC' - t.task_interval + (interval '1 hour' * (1 + floor(random() * 23))) + (interval '1 minute' * floor(random() * 60))
FROM
    users u,
    task_intervals t;
