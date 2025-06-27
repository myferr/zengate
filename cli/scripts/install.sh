#!/bin/bash

set -e

REPO="https://github.com/myferr/zengate"
INSTALL_DIR="$HOME/.zengate"
BIN_PATH="$HOME/.local/bin"

echo "Cloning Zengate..."
git clone "$REPO" "$INSTALL_DIR"

cd "$INSTALL_DIR"
make build

mkdir -p "$BIN_PATH"
cp zengate "$BIN_PATH"

if ! echo "$PATH" | grep -q "$BIN_PATH"; then
  echo "export PATH=\"$BIN_PATH:\$PATH\"" >> "$HOME/.bashrc"
  echo "Added $BIN_PATH to PATH in .bashrc. Restart your terminal or run: source ~/.bashrc"
fi

echo "Zengate installed at $BIN_PATH/zengate"
