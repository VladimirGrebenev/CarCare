#!/usr/bin/env bash
set -euo pipefail

COMPOSE_CMD="${COMPOSE_CMD:-docker-compose}"
PROJECT_NAME="${COMPOSE_PROJECT_NAME:-carcare}"

echo "[carcare] compose v1 safe rebuild start (project: ${PROJECT_NAME})"

${COMPOSE_CMD} down --remove-orphans || true

stale_ids="$(docker ps -aq --filter "label=com.docker.compose.project=${PROJECT_NAME}" || true)"
if [[ -n "${stale_ids}" ]]; then
  echo "[carcare] removing stale compose containers: ${stale_ids}"
  docker rm -f ${stale_ids} >/dev/null || true
fi

${COMPOSE_CMD} rm -f -s -v backend db frontend || true
${COMPOSE_CMD} build --no-cache backend db frontend
${COMPOSE_CMD} up -d --force-recreate backend db frontend

${COMPOSE_CMD} ps
echo "[carcare] compose v1 safe rebuild done"
