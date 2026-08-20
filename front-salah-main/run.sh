#!/usr/bin/env bash
set -e

cleanup() {
  echo "Apagando backend y Apagando frontend"
  kill "${BACKEND_PID:-}" 2>/dev/null || true
  exit
}

trap cleanup SIGINT SIGTERM

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_PORT="${SERVER_PORT:-5000}"
export NUXT_BACKEND_URL="${NUXT_BACKEND_URL:-http://127.0.0.1:${BACKEND_PORT}}"
export GOCACHE="${GOCACHE:-${ROOT_DIR}/.gocache}"

cd "${ROOT_DIR}/backend-salah-main"
echo "Levantando el backend.."
go run . &
BACKEND_PID=$!
echo "Backend en funcionamiento en ${NUXT_BACKEND_URL}"

cd "${ROOT_DIR}/front-salah-main"
echo "Levantando el frontend.."
npm run dev

cleanup
