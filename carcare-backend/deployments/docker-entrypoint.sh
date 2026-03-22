#!/bin/sh
set -e


# Run DB migrations (автоматически, если бинарь есть)
if [ -f /usr/local/bin/migrate ]; then
	/usr/local/bin/migrate -path ./migration -database "$DATABASE_URL" up
fi

exec "$@"
