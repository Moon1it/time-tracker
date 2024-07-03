CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE People (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    passport_serie  INT NOT NULL,
    passport_number INT NOT NULL,
    surname VARCHAR(50),
    name VARCHAR(50),
    patronymic VARCHAR(50),
    address TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc') NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc') NOT NULL,
    UNIQUE (passport_serie, passport_number)
);
