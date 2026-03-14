# Install gen from GitHub Releases. No admin. Copies to user profile.
# Usage: irm https://raw.githubusercontent.com/franco/gen/main/scripts/install.ps1 | iex
# Or: Invoke-WebRequest -Uri <url> -UseBasicParsing | Invoke-Expression

$ErrorActionPreference = "Stop"
$RepoUrl = if ($env:GEN_REPO_URL) { $env:GEN_REPO_URL } else { "https://github.com/franco/gen" }
$BaseUrl = "$RepoUrl/releases/latest/download"
$ArchiveName = "gen_windows_amd64.zip"
$BinaryName = "gen.exe"

$arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE", "Machine")
if ($arch -notmatch "AMD64|x64") {
    Write-Error "gen: apenas Windows amd64 é suportado (arquitetura atual: $arch)"
    exit 1
}

$url = "$BaseUrl/$ArchiveName"
$tmpDir = Join-Path $env:TEMP "gen-install-$(Get-Random)"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
try {
    Write-Host "Baixando $url ..."
    Invoke-WebRequest -Uri $url -UseBasicParsing -OutFile (Join-Path $tmpDir $ArchiveName)
    Expand-Archive -Path (Join-Path $tmpDir $ArchiveName) -DestinationPath $tmpDir -Force
    $binDir = Join-Path $env:USERPROFILE "bin"
    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }
    $dest = Join-Path $binDir $BinaryName
    Copy-Item (Join-Path $tmpDir $BinaryName) -Destination $dest -Force
    Write-Host "Instalado em $dest"
    $pathEnv = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($pathEnv -notlike "*$binDir*") {
        Write-Host ""
        Write-Host "O diretório $binDir não está no PATH. Adicione em Configurações > Sistema > Sobre > Configurações avançadas > Variáveis de ambiente."
        Write-Host "Ou execute: [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$binDir', 'User')"
        Write-Host ""
    }
} finally {
    Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
}
