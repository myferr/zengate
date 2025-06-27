#!/bin/bash

set -e

REPO_URL="https://github.com/myferr/zengate"
INSTALL_DIR="$HOME/.zengate"
BIN_DIR="$HOME/.local/bin"

echo "Cloning Zengate into $INSTALL_DIR..."
git clone --depth 1 "$REPO_URL" "$INSTALL_DIR"

cd "$INSTALL_DIR/cli"
echo "Building CLI..."
go build -o zengate .

mkdir -p "$BIN_DIR"
cp zengate "$BIN_DIR"

# Add to PATH if not already
if ! echo "$PATH" | grep -q "$BIN_DIR"; then
  SHELL_RC="$HOME/.bashrc"
  [ -n "$ZSH_VERSION" ] && SHELL_RC="$HOME/.zshrc"
  echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$SHELL_RC"
  echo "Added $BIN_DIR to PATH in $SHELL_RC"
fi

echo "Zengate CLI installed at $BIN_DIR/zengate"
echo "Run: zengate help"
