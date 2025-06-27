$BIN_NAME = "zengate.exe"
$INSTALL_DIR = "$env:USERPROFILE\bin"
$GO_BUILD_OUTPUT = ".\$BIN_NAME"

Write-Host "Building $BIN_NAME..."
go build -o $GO_BUILD_OUTPUT
if ($LASTEXITCODE -ne 0) {
  Write-Error "Go build failed! Make sure Go is installed."
  exit 1
}

if (-Not (Test-Path $INSTALL_DIR)) {
  Write-Host "Creating directory $INSTALL_DIR"
  New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

Write-Host "Installing $BIN_NAME to $INSTALL_DIR"
Move-Item -Force $GO_BUILD_OUTPUT "$INSTALL_DIR\$BIN_NAME"

if (-Not ($env:PATH.Split(';') -contains $INSTALL_DIR)) {
  Write-Host ""
  Write-Warning "$INSTALL_DIR is not in your PATH environment variable."
  Write-Host "You may want to add it. You can do this via:"
  Write-Host "  [Environment]::SetEnvironmentVariable('PATH', $env:PATH + ';$INSTALL_DIR', 'User')"
  Write-Host "Then restart your terminal or computer."
}

Write-Host "$BIN_NAME installed successfully!"
Write-Host "You can now run it with: $BIN_NAME"
