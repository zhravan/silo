# Install silo from GitHub releases (Windows).
# Usage: irm https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.ps1 | iex
# Or: $env:SILO_VERSION="0.0.1"; irm ... | iex

$ErrorActionPreference = "Stop"
$Repo = "zhravan/silo"
$ApiBase = "https://api.github.com/repos/$Repo"

function Get-LatestVersion {
  $r = Invoke-RestMethod -Uri "$ApiBase/releases/latest" -MaximumRetentionInSeconds 10
  $tag = $r.tag_name
  if ($tag -match "^v?(.+)$") { return $Matches[1] }
  return $tag.TrimStart("v")
}

$Version = if ($env:SILO_VERSION) { $env:SILO_VERSION.Trim() } else { Get-LatestVersion }
if (-not $Version) {
  Write-Error "Could not determine version. Set SILO_VERSION or ensure $Repo has a release."
  exit 1
}

$Tag = if ($Version -notmatch "^v") { "v$Version" } else { $Version }
$BaseUrl = "https://github.com/$Repo/releases/download/$Tag"

$Arch = if ([Environment]::Is64BitOperatingSystem) {
  $proc = (Get-CimInstance Win32_Processor).Architecture
  if ($proc -eq 12) { "arm64" } else { "amd64" }
} else {
  "amd64"
}

$Binary = "silo-windows-$Arch.exe"
$Url = "$BaseUrl/$Binary"
Write-Host "Installing silo $Tag ($Arch) from $Url"

$InstallDir = if ($env:SILO_PREFIX) { $env:SILO_PREFIX } else { "$env:LOCALAPPDATA\Programs\silo" }
$Dest = Join-Path $InstallDir "silo.exe"

$null = New-Item -ItemType Directory -Force -Path $InstallDir
Invoke-WebRequest -Uri $Url -OutFile $Dest -UseBasicParsing

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$InstallDir*") {
  [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
  $env:Path = "$env:Path;$InstallDir"
  Write-Host "Added $InstallDir to user PATH."
}

Write-Host "Installed to $Dest"
& $Dest --help 2>$null
