# Install gen from GitHub Releases. No admin. Copies to user profile.
# Usage: irm https://raw.githubusercontent.com/francofabio/gen/main/scripts/install.ps1 | iex
# Or: Invoke-WebRequest -Uri <url> -UseBasicParsing | Invoke-Expression

$ErrorActionPreference = "Stop"
$RepoUrl = if ($env:GEN_REPO_URL) { $env:GEN_REPO_URL } else { "https://github.com/francofabio/gen" }
$BaseUrl = "$RepoUrl/releases/latest/download"
$ArchiveName = "gen_windows_amd64.zip"
$BinaryName = "gen.exe"

$arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE", "Machine")
if ($arch -notmatch "AMD64|x64") {
    Write-Error "gen: only Windows amd64 is supported (current architecture: $arch)"
    exit 1
}

$url = "$BaseUrl/$ArchiveName"
$tmpDir = Join-Path $env:TEMP "gen-install-$(Get-Random)"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
try {
    Write-Host "Downloading $url ..."
    Invoke-WebRequest -Uri $url -UseBasicParsing -OutFile (Join-Path $tmpDir $ArchiveName)
    Expand-Archive -Path (Join-Path $tmpDir $ArchiveName) -DestinationPath $tmpDir -Force
    $binDir = Join-Path $env:USERPROFILE "bin"
    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }
    $dest = Join-Path $binDir $BinaryName
    Copy-Item (Join-Path $tmpDir $BinaryName) -Destination $dest -Force
    Write-Host "Installed at $dest"
    $pathEnv = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($pathEnv -notlike "*$binDir*") {
        Write-Host ""
        Write-Host "The directory $binDir is not in your PATH. Add it via Settings > System > About > Advanced system settings > Environment variables."
        Write-Host "Or run: [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$binDir', 'User')"
        Write-Host ""
    }
} finally {
    Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
}
