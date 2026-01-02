<#
.SYNOPSIS
    Install migrate-tool from GitHub releases

.DESCRIPTION
    Downloads and installs migrate-tool from GitHub releases with automatic
    architecture detection, checksum verification, and optional PATH addition.

.PARAMETER Version
    Specific version to install (e.g., v0.0.4). If not specified, installs latest.

.PARAMETER InstallDir
    Installation directory. Default: $env:LOCALAPPDATA\migrate-tool

.PARAMETER AddToPath
    Add installation directory to user PATH environment variable.

.EXAMPLE
    .\install.ps1
    Install latest version to default location

.EXAMPLE
    .\install.ps1 -Version v0.0.4
    Install specific version

.EXAMPLE
    .\install.ps1 -AddToPath
    Install and add to PATH
#>

[CmdletBinding()]
param(
    [string]$Version = "",
    [string]$InstallDir = "$env:LOCALAPPDATA\migrate-tool",
    [switch]$AddToPath,
    [switch]$Help
)

# Configuration
$Script:Repo = "cesc1802/migration-tool"
$Script:BinaryName = "migrate-tool"
$Script:GitHubApi = "https://api.github.com"
$Script:GitHubDownload = "https://github.com"

# Logging functions
function Write-Info { param([string]$Message) Write-Host "[INFO] $Message" -ForegroundColor Blue }
function Write-Ok { param([string]$Message) Write-Host "[OK] $Message" -ForegroundColor Green }
function Write-Warn { param([string]$Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Err { param([string]$Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }

function Exit-WithError {
    param([string]$Message)
    Write-Err $Message
    exit 1
}

function Show-Usage {
    @"
Install migrate-tool from GitHub releases

Usage: .\install.ps1 [OPTIONS]

Options:
    -Version VERSION    Install specific version (e.g., v0.0.4)
    -InstallDir DIR     Install directory (default: $env:LOCALAPPDATA\migrate-tool)
    -AddToPath          Add install directory to user PATH
    -Help               Show this help message

Examples:
    .\install.ps1                       Install latest version
    .\install.ps1 -Version v0.0.4       Install specific version
    .\install.ps1 -AddToPath            Install and add to PATH
    .\install.ps1 -InstallDir C:\Tools  Install to custom directory
"@
    exit 0
}

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Exit-WithError "Unsupported architecture: $arch (supported: AMD64, ARM64)" }
    }
}

function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    try {
        $response = Invoke-RestMethod -Uri "$Script:GitHubApi/repos/$Script:Repo/releases/latest" -ErrorAction Stop
        $version = $response.tag_name
        if (-not $version) {
            Exit-WithError "Failed to parse version from GitHub API response"
        }
        if ($version -notmatch "^v\d") {
            Exit-WithError "Invalid version format from API: $version"
        }
        return $version
    }
    catch {
        Exit-WithError "Failed to fetch latest release. Check network or try: -Version vX.Y.Z"
    }
}

function Get-FileFromUrl {
    param(
        [string]$Url,
        [string]$OutFile
    )
    try {
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $Url -OutFile $OutFile -ErrorAction Stop
        $ProgressPreference = 'Continue'
    }
    catch {
        Exit-WithError "Failed to download from $Url"
    }
}

function Test-Checksum {
    param(
        [string]$ArchivePath,
        [string]$ChecksumsPath
    )

    $archiveName = Split-Path $ArchivePath -Leaf
    $checksumContent = Get-Content $ChecksumsPath

    # Find expected checksum for this file
    $expectedLine = $checksumContent | Where-Object { $_ -match $archiveName }
    if (-not $expectedLine) {
        Exit-WithError "Checksum not found for $archiveName in checksums.txt"
    }
    $expected = ($expectedLine -split '\s+')[0].ToLower()

    # Calculate actual checksum
    $actual = (Get-FileHash -Path $ArchivePath -Algorithm SHA256).Hash.ToLower()

    if ($expected -ne $actual) {
        Write-Err "Checksum verification failed!"
        Write-Err "Expected: $expected"
        Write-Err "Actual:   $actual"
        Exit-WithError "Downloaded file may be corrupted or tampered with"
    }
    Write-Ok "Checksum verified"
}

function Expand-TarGz {
    param(
        [string]$ArchivePath,
        [string]$DestinationPath
    )

    # Windows 10 1803+ has tar built-in
    $tarPath = Get-Command tar -ErrorAction SilentlyContinue
    if ($tarPath) {
        try {
            & tar -xzf $ArchivePath -C $DestinationPath 2>&1 | Out-Null
            if ($LASTEXITCODE -ne 0) {
                throw "tar extraction failed"
            }
            return
        }
        catch {
            Write-Warn "tar extraction failed, trying .NET fallback..."
        }
    }

    # Fallback: Use .NET for older Windows
    # First decompress gzip, then extract tar
    $gzipStream = $null
    $tarStream = $null
    $decompressionStream = $null
    $tarPath = Join-Path $DestinationPath "archive.tar"

    try {
        Add-Type -AssemblyName System.IO.Compression.FileSystem

        $gzipStream = [System.IO.File]::OpenRead($ArchivePath)
        $tarStream = [System.IO.File]::Create($tarPath)
        $decompressionStream = New-Object System.IO.Compression.GZipStream($gzipStream, [System.IO.Compression.CompressionMode]::Decompress)
        $decompressionStream.CopyTo($tarStream)
    }
    finally {
        if ($decompressionStream) { $decompressionStream.Dispose() }
        if ($tarStream) { $tarStream.Dispose() }
        if ($gzipStream) { $gzipStream.Dispose() }
    }

    try {
        # Simple tar extraction (handles basic tar format)
        $tarBytes = [System.IO.File]::ReadAllBytes($tarPath)
        $offset = 0
        while ($offset -lt $tarBytes.Length) {
            # Read header (512 bytes)
            if ($offset + 512 -gt $tarBytes.Length) { break }

            # Check for empty block (end of archive)
            $isEmpty = $true
            for ($i = 0; $i -lt 512; $i++) {
                if ($tarBytes[$offset + $i] -ne 0) { $isEmpty = $false; break }
            }
            if ($isEmpty) { break }

            # Get filename (first 100 bytes, null-terminated)
            $nameBytes = $tarBytes[$offset..($offset + 99)]
            $nameEnd = [Array]::IndexOf($nameBytes, [byte]0)
            if ($nameEnd -lt 0) { $nameEnd = 100 }
            $fileName = [System.Text.Encoding]::ASCII.GetString($nameBytes, 0, $nameEnd).Trim()

            # Get file size (octal, bytes 124-135) with validation
            $sizeStr = [System.Text.Encoding]::ASCII.GetString($tarBytes, $offset + 124, 11).Trim()
            if ($sizeStr -notmatch '^[0-7]*$') {
                throw "Invalid tar file size format"
            }
            $fileSize = if ($sizeStr) { [Convert]::ToInt64($sizeStr, 8) } else { 0 }
            if ($fileSize -lt 0 -or $fileSize -gt 100MB) {
                throw "Invalid tar file size: $fileSize"
            }

            # Get file type (byte 156)
            $fileType = $tarBytes[$offset + 156]

            $offset += 512  # Move past header

            # Extract regular files (type '0' or null)
            if (($fileType -eq 48 -or $fileType -eq 0) -and $fileSize -gt 0 -and $fileName) {
                $outPath = Join-Path $DestinationPath $fileName
                $outDir = Split-Path $outPath -Parent
                if (-not (Test-Path $outDir)) {
                    New-Item -ItemType Directory -Path $outDir -Force | Out-Null
                }
                $fileBytes = $tarBytes[$offset..($offset + $fileSize - 1)]
                [System.IO.File]::WriteAllBytes($outPath, $fileBytes)
            }

            # Move to next entry (512-byte aligned)
            $offset += [Math]::Ceiling($fileSize / 512) * 512
        }
    }
    finally {
        Remove-Item $tarPath -Force -ErrorAction SilentlyContinue
    }
}

function Add-ToUserPath {
    param([string]$Directory)

    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $normalizedDir = $Directory.ToLower().TrimEnd('\')
    $pathEntries = $currentPath -split ';' | ForEach-Object { $_.ToLower().TrimEnd('\') }
    if ($pathEntries -contains $normalizedDir) {
        Write-Info "$Directory is already in PATH"
        return
    }

    try {
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Ok "Added $Directory to user PATH"
        Write-Warn "Restart your terminal for PATH changes to take effect"
    }
    catch {
        Write-Warn "Failed to add to PATH: $_"
        Write-Info "Manually add to PATH: $Directory"
    }
}

function Main {
    if ($Help) {
        Show-Usage
    }

    Write-Info "Installing migrate-tool"
    Write-Host ""

    # Detect architecture
    $arch = Get-Architecture
    Write-Info "Detected architecture: windows/$arch"

    # Get version
    if (-not $Version) {
        $Version = Get-LatestVersion
    }
    $versionNum = $Version -replace '^v', ''
    Write-Info "Version: $Version"

    # Create temp directory
    $tempDir = Join-Path $env:TEMP "migrate-tool-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    try {
        # Build download URLs
        $archiveName = "${Script:BinaryName}_${versionNum}_windows_${arch}.tar.gz"
        $archiveUrl = "$Script:GitHubDownload/$Script:Repo/releases/download/$Version/$archiveName"
        $checksumsUrl = "$Script:GitHubDownload/$Script:Repo/releases/download/$Version/checksums.txt"

        # Download archive and checksums
        Write-Info "Downloading $archiveName..."
        Get-FileFromUrl -Url $archiveUrl -OutFile (Join-Path $tempDir $archiveName)
        Write-Ok "Downloaded archive"

        Write-Info "Downloading checksums..."
        Get-FileFromUrl -Url $checksumsUrl -OutFile (Join-Path $tempDir "checksums.txt")
        Write-Ok "Downloaded checksums"

        # Verify checksum
        Write-Info "Verifying checksum..."
        Test-Checksum -ArchivePath (Join-Path $tempDir $archiveName) -ChecksumsPath (Join-Path $tempDir "checksums.txt")

        # Extract archive
        Write-Info "Extracting archive..."
        Expand-TarGz -ArchivePath (Join-Path $tempDir $archiveName) -DestinationPath $tempDir
        Write-Ok "Extracted archive"

        # Create install directory
        if (-not (Test-Path $InstallDir)) {
            try {
                New-Item -ItemType Directory -Path $InstallDir -Force -ErrorAction Stop | Out-Null
            }
            catch {
                Exit-WithError "Failed to create install directory: $InstallDir"
            }
        }

        # Install binary
        Write-Info "Installing to $InstallDir..."
        $binaryExt = ".exe"
        $sourceBinary = Join-Path $tempDir "${Script:BinaryName}$binaryExt"
        if (-not (Test-Path $sourceBinary)) {
            # Try without extension (some archives don't include .exe)
            $sourceBinary = Join-Path $tempDir $Script:BinaryName
        }
        if (-not (Test-Path $sourceBinary)) {
            Exit-WithError "Binary not found in archive"
        }

        $destBinary = Join-Path $InstallDir "${Script:BinaryName}$binaryExt"
        Copy-Item $sourceBinary $destBinary -Force
        Write-Ok "Installed ${Script:BinaryName} to $InstallDir"

        # Add to PATH if requested
        if ($AddToPath) {
            Add-ToUserPath -Directory $InstallDir
        }

        Write-Host ""
        Write-Ok "Installation complete!"
        Write-Host ""

        # Verify installation
        if (Get-Command $Script:BinaryName -ErrorAction SilentlyContinue) {
            $ver = & $Script:BinaryName version 2>&1
            Write-Info "Version: $ver"
        }
        else {
            if (-not $AddToPath) {
                Write-Warn "$InstallDir is not in your PATH"
                Write-Info "Run with -AddToPath or manually add: $InstallDir"
            }
        }
    }
    finally {
        # Cleanup temp directory
        if (Test-Path $tempDir) {
            Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

Main
