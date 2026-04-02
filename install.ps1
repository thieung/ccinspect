# ccinspect installer for Windows
# Usage: irm https://raw.githubusercontent.com/thieung/ccinspect/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "thieung/ccinspect"
$Binary = "ccinspect"
$InstallDir = "$env:LOCALAPPDATA\ccinspect"

function Write-Info($msg)  { Write-Host "▸ $msg" -ForegroundColor Cyan }
function Write-Ok($msg)    { Write-Host "✓ $msg" -ForegroundColor Green }
function Write-Warn($msg)  { Write-Host "⚠ $msg" -ForegroundColor Yellow }
function Write-Err($msg)   { Write-Host "✗ $msg" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "╭─────────────────────────────────╮" -ForegroundColor Cyan
Write-Host "│   ccinspect installer            │" -ForegroundColor Cyan
Write-Host "╰─────────────────────────────────╯" -ForegroundColor Cyan
Write-Host ""

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Err "32-bit systems are not supported."
}
Write-Info "Detected: windows/$Arch"

# Get latest version
Write-Info "Fetching latest version..."
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
$Version = $Release.tag_name -replace '^v', ''

if (-not $Version) {
    Write-Err "Could not determine latest version."
}
Write-Ok "Version: v$Version"

# Download
$Archive = "${Binary}_${Version}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/v$Version/$Archive"
$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "ccinspect-install-$(Get-Random)")

Write-Info "Downloading $Archive..."
try {
    Invoke-WebRequest -Uri $Url -OutFile (Join-Path $TmpDir $Archive) -UseBasicParsing
} catch {
    Write-Err "Download failed. Check https://github.com/$Repo/releases"
}
Write-Ok "Downloaded"

# Extract
Write-Info "Extracting..."
Expand-Archive -Path (Join-Path $TmpDir $Archive) -DestinationPath $TmpDir -Force
Write-Ok "Extracted"

# Install
Write-Info "Installing to $InstallDir..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
Copy-Item -Path (Join-Path $TmpDir "$Binary.exe") -Destination (Join-Path $InstallDir "$Binary.exe") -Force
Write-Ok "Installed to $InstallDir\$Binary.exe"

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Info "Adding $InstallDir to PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path = "$env:Path;$InstallDir"
    Write-Ok "Added to PATH (restart terminal for full effect)"
}

# Cleanup
Remove-Item -Recurse -Force $TmpDir

# Verify
Write-Host ""
Write-Ok (& (Join-Path $InstallDir "$Binary.exe") --version 2>&1)
Write-Host ""
Write-Host "🎉 Installation complete!" -ForegroundColor Green
Write-Host "   Run 'ccinspect --help' to get started." -ForegroundColor Cyan
Write-Host "   (You may need to restart your terminal)" -ForegroundColor Yellow
Write-Host ""
