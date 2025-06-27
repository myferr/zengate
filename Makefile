CLI_DIR := cli
DESKTOP_DIR := desktop
CLI_BINARY := zengate
BIN_DIR := $(HOME)/.local/bin

.PHONY: all cli desktop install clean

all: cli desktop

## Build the CLI
cli:
	@echo "Building CLI..."
	cd $(CLI_DIR) && go build -o $(CLI_BINARY)

## Build the Tauri desktop app
desktop:
	@echo " Building desktop app..."
	cd $(DESKTOP_DIR) && bun install && bunx tauri build

## Install CLI to ~/.local/bin
install: cli
	@mkdir -p $(BIN_DIR)
	@cp $(CLI_DIR)/$(CLI_BINARY) $(BIN_DIR)/
	@echo "✅ Installed CLI to $(BIN_DIR)"

## Clean CLI binary
clean:
	@echo "Cleaning up CLI binary..."
	@rm -f $(CLI_DIR)/$(CLI_BINARY)
