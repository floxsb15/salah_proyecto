#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
ENV_FILE="$SCRIPT_DIR/.env"
BASE_COMPOSE="$SCRIPT_DIR/compose.yml"
TEST_COMPOSE="$SCRIPT_DIR/compose.local-test.yml"

if [ ! -f "$ENV_FILE" ]; then
  echo "Falta $ENV_FILE. Ejecute primero ./deploy/prepare-vps-test.sh" >&2
  exit 1
fi

chmod 600 "$ENV_FILE"
cd "$REPO_DIR"

docker compose --env-file "$ENV_FILE" \
  -f "$BASE_COMPOSE" \
  -f "$TEST_COMPOSE" \
  config --quiet

docker compose --env-file "$ENV_FILE" \
  -f "$BASE_COMPOSE" \
  -f "$TEST_COMPOSE" \
  up -d --build --wait --wait-timeout 300 postgres backend frontend

docker compose --env-file "$ENV_FILE" \
  -f "$BASE_COMPOSE" \
  -f "$TEST_COMPOSE" \
  ps postgres backend frontend

echo "Aplicacion lista en el puerto definido por FRONTEND_HOST_PORT (3100 por defecto)."
echo "El puerto solo escucha en 127.0.0.1 del VPS. Caddy no fue iniciado."
echo "Consulte la clave inicial localmente con: grep '^INITIAL_ADMIN_PASSWORD=' deploy/.env"
