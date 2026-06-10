-- Миграция: добавляем поле bill_number в таблицу fines
-- Выполнить на существующей БД: psql -U postgres -d carcare -f migration_002_add_bill_number.sql

ALTER TABLE fines ADD COLUMN IF NOT EXISTS bill_number VARCHAR(50) NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_fines_bill_number ON fines(bill_number);