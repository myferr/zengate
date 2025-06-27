#!/usr/bin/env bash

set -e

BIN_NAME="zengate"
INSTALL_DIR="$HOME/.local/bin"
GO_BUILD_OUTPUT="./$BIN_NAME"

echo "Building $BIN_NAME..."
go build -o "$GO_BUILD_OUTPUT" || { echo "Go build failed! Make sure Go is installed."; exit 1; }

mkdir -p "$INSTALL_DIR"

echo "Installing $BIN_NAME to $INSTALL_DIR"
mv -f "$GO_BUILD_OUTPUT" "$INSTALL_DIR/"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo ""
  echo "⚠️  Warning: $INSTALL_DIR is not in your PATH."
  echo "Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
  echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "$BIN_NAME installed successfully!"
echo "You can now run it with: $BIN_NAME"
