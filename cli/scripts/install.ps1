$Repo = "https://github.com/myferr/zengate"
$InstallDir = "$env:USERPROFILE\.zengate"
$BinPath = "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps"

Write-Host "Cloning Zengate..."
git clone $Repo $InstallDir

Set-Location $InstallDir
& "go" "build" "-o" "zengate.exe"

Copy-Item -Path ".\zengate.exe" -Destination $BinPath -Force

if (-not ($env:PATH -like "*$BinPath*")) {
    Write-Host "`n⚠️ You may need to add $BinPath to your PATH manually."
} else {
    Write-Host "Zengate installed at $BinPath\zengate.exe"
}
