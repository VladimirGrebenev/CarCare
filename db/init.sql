-- CarCare database initialization
-- Выполняется автоматически PostgreSQL при первом запуске контейнера
-- (docker-entrypoint-initdb.d)

CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY,
    email       VARCHAR(255) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    role        VARCHAR(50)  NOT NULL,
    password_hash VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS cars (
    id      UUID PRIMARY KEY,
    user_id UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    brand   VARCHAR(100) NOT NULL,
    model   VARCHAR(100) NOT NULL,
    year    INT          NOT NULL,
    vin     VARCHAR(100) NOT NULL UNIQUE,
    plate   VARCHAR(20)
);
CREATE INDEX IF NOT EXISTS idx_cars_user_id ON cars(user_id);

CREATE TABLE IF NOT EXISTS fuel_events (
    id     UUID PRIMARY KEY,
    car_id UUID           NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    volume NUMERIC(10,2)  NOT NULL,
    price  NUMERIC(10,2)  NOT NULL,
    type   VARCHAR(50)    NOT NULL,
    date   DATE           NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_fuel_events_car_id ON fuel_events(car_id);

CREATE TABLE IF NOT EXISTS maintenance_events (
    id          UUID PRIMARY KEY,
    car_id      UUID          NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    type        VARCHAR(100)  NOT NULL,
    date        DATE          NOT NULL,
    cost        NUMERIC(10,2) NOT NULL,
    description TEXT          NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_maintenance_events_car_id ON maintenance_events(car_id);

CREATE TABLE IF NOT EXISTS fines (
    id          UUID PRIMARY KEY,
    car_id      UUID          NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    amount      NUMERIC(10,2) NOT NULL,
    type        VARCHAR(100)  NOT NULL,
    date        DATE          NOT NULL,
    status      VARCHAR(20)   NOT NULL DEFAULT 'unpaid',
    description TEXT          NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_fines_car_id ON fines(car_id);

CREATE TABLE IF NOT EXISTS reports (
    id         SERIAL    PRIMARY KEY,
    type       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    data       JSONB        NOT NULL
);
