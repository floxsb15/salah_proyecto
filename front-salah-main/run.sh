#! /bin/bash
set -e

cleanup() {
  echo "Apagando backend y Apagando frontend"
  kill $BACKEND_PID 2>/dev/null || true
  exit
}

trap cleanup SIGINT SIGTERM

cd backend
echo "Levantando el backend.."
./backend_restaurant &
BACKEND_PID=$!
echo "Backend en funcionamiento"
cd ..
echo "Levantando el frontend.."
npm run dev

cleanup