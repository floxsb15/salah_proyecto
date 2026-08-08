#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ENV_FILE="$SCRIPT_DIR/.env"
ADMIN_USER=${1:-administrador}

if [ -e "$ENV_FILE" ]; then
  echo "No se reemplazo $ENV_FILE porque ya existe." >&2
  echo "Si desea regenerarlo, mueva primero el archivo existente a una ubicacion segura." >&2
  exit 1
fi

for command_name in docker openssl; do
  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Falta el comando requerido: $command_name" >&2
    exit 1
  fi
done

docker compose version >/dev/null

AUTH_TOKEN_SECRET=$(openssl rand -hex 48)
DB_PASSWORD=$(openssl rand -hex 32)
INITIAL_ADMIN_PASSWORD=$(openssl rand -hex 20)

umask 077
{
  printf '%s\n' 'DOMAIN=localhost'
  printf 'AUTH_TOKEN_SECRET=%s\n' "$AUTH_TOKEN_SECRET"
  printf '%s\n' ''
  printf '%s\n' 'DB_HOST=postgres'
  printf '%s\n' 'DB_PORT=5432'
  printf '%s\n' 'DB_USER=salah'
  printf 'DB_PASSWORD=%s\n' "$DB_PASSWORD"
  printf '%s\n' 'DB_NAME=salah'
  printf '%s\n' 'DB_SSLMODE=disable'
  printf '%s\n' ''
  printf '%s\n' 'RUN_MIGRATIONS=true'
  printf 'INITIAL_ADMIN_USER=%s\n' "$ADMIN_USER"
  printf 'INITIAL_ADMIN_PASSWORD=%s\n' "$INITIAL_ADMIN_PASSWORD"
} > "$ENV_FILE"
chmod 600 "$ENV_FILE"

echo "Entorno de prueba creado en $ENV_FILE con permisos 600."
echo "Los secretos no se imprimieron. Ejecute ahora: ./deploy/start-vps-test.sh"
