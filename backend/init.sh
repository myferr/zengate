#!/bin/bash
set -e

echo "[+] Checking Bun..."
if ! command -v bun &>/dev/null; then
  echo "[-] Bun not found, install it from https://bun.sh"
  exit 1
fi

echo "[+] Installing Bun dependencies..."
cd backend
bun install
cd ..

echo "[+] Checking Go..."
if ! command -v go &>/dev/null; then
  echo "[-] Go not found, please install Go 1.21+"
  exit 1
fi

echo "[+] Installing Go tunnel dependencies..."
cd backend
go mod tidy

cd ../

echo "[✓] Initialization complete."
