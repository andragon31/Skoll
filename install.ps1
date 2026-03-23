$ErrorActionPreference = "Stop"

$REPO = "andragon31/skoll"
$BIN = "skoll-windows-amd64.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\Programs\skoll"
$EXE_PATH = "$INSTALL_DIR\skoll.exe"

$VERSION = "v0.1.0"

Clear-Host
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  Skoll Installer $VERSION" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  >> RSAW Orchestration Layer <<" -ForegroundColor Gray
Write-Host "  -----------------------------" -ForegroundColor Gray
Write-Host ""

if (Test-Path $EXE_PATH) {
    Write-Host "Skoll already installed. Updating..." -ForegroundColor Yellow
}

Write-Host "[1/3] Compiling Skoll (requires Go)..." -ForegroundColor DarkCyan
try {
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
    }
    go build -o $EXE_PATH ./cmd/skoll
} catch {
    Write-Host "Error compiling: $_" -ForegroundColor Red
    exit 1
}

Write-Host "[2/3] Installing to $INSTALL_DIR..." -ForegroundColor DarkCyan

Write-Host "[3/3] Adding to PATH..." -ForegroundColor DarkCyan

$currentMachinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$currentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")

$alreadyInMachine = $currentMachinePath -split ";" | Where-Object { $_.Trim() -eq $INSTALL_DIR }
$alreadyInUser = $currentUserPath -split ";" | Where-Object { $_.Trim() -eq $INSTALL_DIR }

$pathAdded = $false

if (-not $alreadyInMachine) {
    try {
        [Environment]::SetEnvironmentVariable("Path", "$INSTALL_DIR;$currentMachinePath", "Machine")
        Write-Host "  Added to System PATH (Machine)" -ForegroundColor Green
        $pathAdded = $true
    } catch {
        Write-Host "  No admin rights - using User PATH" -ForegroundColor Yellow
    }
}

if (-not $alreadyInUser) {
    [Environment]::SetEnvironmentVariable("Path", "$INSTALL_DIR;$currentUserPath", "User")
    Write-Host "  Added to User PATH" -ForegroundColor Green
    $pathAdded = $true
}

$env:Path = "$INSTALL_DIR;$currentMachinePath;$currentUserPath"

Write-Host ""
Write-Host "  [ Verification ]" -ForegroundColor DarkCyan
Write-Host "  ----------------" -ForegroundColor Gray
Write-Host ""

try {
    & $EXE_PATH version
} catch {
    Write-Host "Version check failed: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Green
Write-Host "  skoll init            # Initialize RSAW in project"
Write-Host "  skoll setup cursor    # Setup for Cursor"
Write-Host "  skoll tui             # Open Dashboard"
Write-Host ""

if ($pathAdded) {
    Write-Host "NOTE: If 'skoll' is not found, open a new PowerShell window." -ForegroundColor Yellow
}
