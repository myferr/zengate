$Repo = "https://github.com/myferr/zengate"
$InstallDir = "$env:USERPROFILE\.zengate"
$BinDir = "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps"

Write-Host "Cloning Zengate into $InstallDir..."
git clone --depth 1 $Repo $InstallDir

Set-Location "$InstallDir\cli"
Write-Host "Building zengate.exe..."
go build -o zengate.exe .

Copy-Item -Path ".\zengate.exe" -Destination $BinDir -Force

Write-Host "Zengate CLI installed at $BinDir\zengate.exe"
Write-Host "You may need to restart your terminal if PATH hasn't updated"
