#!/bin/sh
set -e

# Схема БД создаётся через db/init.sql (docker-entrypoint-initdb.d PostgreSQL)
exec "$@"
