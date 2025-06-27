#!/bin/bash
set -e

echo "[+] Loading env vars from backend/.env"
export $(grep -v '^#' backend/.env | xargs)

echo "[+] Starting ngrok tunnel and backend..."

cd backend/ngrok
go run main.go &

NGROK_PID=$!

cd ../

# Wait a few seconds for ngrok to establish tunnel before starting backend
sleep 5

bun run src/index.ts

kill $NGROK_PID
