-- Миграция: добавляем поля sts и bill_number
-- Выполнить на существующей БД: psql -U postgres -d carcare -f migration_001_add_sts.sql

ALTER TABLE cars ADD COLUMN IF NOT EXISTS sts VARCHAR(20);
ALTER TABLE fines ADD COLUMN IF NOT EXISTS bill_number VARCHAR(50) NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_fines_bill_number ON fines(bill_number);