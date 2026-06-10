-- Миграция: добавляем поле sts в таблицу cars
-- Выполнить на существующей БД: psql -U postgres -d carcare -f migration_001_add_sts.sql

ALTER TABLE cars ADD COLUMN IF NOT EXISTS sts VARCHAR(20);