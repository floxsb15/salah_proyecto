#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
ENV_FILE="$SCRIPT_DIR/.env"
BASE_COMPOSE="$SCRIPT_DIR/compose.yml"
TEST_COMPOSE="$SCRIPT_DIR/compose.local-test.yml"

if [ "${1:-}" != "--password-saved" ]; then
  echo "Guarde primero la clave del administrador y ejecute:" >&2
  echo "./deploy/finish-vps-test-setup.sh --password-saved" >&2
  exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
  echo "Falta $ENV_FILE." >&2
  exit 1
fi

umask 077
TEMP_FILE=$(mktemp "$SCRIPT_DIR/.env.XXXXXX")
trap 'rm -f "$TEMP_FILE"' EXIT HUP INT TERM

awk '
  /^RUN_MIGRATIONS=/ { print "RUN_MIGRATIONS=false"; next }
  /^INITIAL_ADMIN_PASSWORD=/ { print "INITIAL_ADMIN_PASSWORD="; next }
  { print }
' "$ENV_FILE" > "$TEMP_FILE"

chmod 600 "$TEMP_FILE"
mv "$TEMP_FILE" "$ENV_FILE"
trap - EXIT HUP INT TERM

cd "$REPO_DIR"
docker compose --env-file "$ENV_FILE" \
  -f "$BASE_COMPOSE" \
  -f "$TEST_COMPOSE" \
  up -d --force-recreate --wait --wait-timeout 120 backend frontend

echo "Migraciones desactivadas y clave inicial retirada del entorno de los contenedores."
