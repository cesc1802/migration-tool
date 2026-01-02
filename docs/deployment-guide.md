# Deployment & Installation Guide

## Unix Installation Script

**Phase 01 Completion:** Automated installation script for Unix-like systems (Linux, macOS).

### Overview

The `scripts/install.sh` script provides a robust, cross-platform method to download and install `migrate-tool` from GitHub releases. It includes:
- Automatic OS/architecture detection (Linux/macOS, amd64/arm64)
- SHA256 checksum verification
- Dependency checking (curl/wget, tar, sha256sum/shasum)
- Automatic fallback between curl and wget
- Colored output with NO_COLOR support
- Installation to custom directories
- Version-specific or latest version installation

### Usage

#### Basic Installation (Latest Version)

```bash
./scripts/install.sh
```

Installs latest `migrate-tool` release to `/usr/local/bin`.

#### System-Wide Installation (Requires sudo)

```bash
sudo ./scripts/install.sh
```

#### User-Local Installation (No sudo required)

```bash
INSTALL_DIR=~/.local/bin ./scripts/install.sh
```

Ensure `~/.local/bin` is in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

#### Install Specific Version

```bash
./scripts/install.sh --version v0.0.3
```

#### Custom Installation Directory

```bash
./scripts/install.sh --dir /opt/migrate-tool
```

### Options

| Option | Description | Example |
|--------|-------------|---------|
| `-v, --version VERSION` | Install specific version | `./install.sh --version v0.0.3` |
| `-d, --dir DIRECTORY` | Custom install directory | `./install.sh --dir ~/.local/bin` |
| `-h, --help` | Show help message | `./install.sh --help` |

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `INSTALL_DIR` | Installation directory | `/usr/local/bin` |
| `NO_COLOR` | Disable colored output | (not set = colored) |

### Examples

```bash
# Install latest to system directory
sudo ./scripts/install.sh

# Install specific version to home directory
INSTALL_DIR=~/.local/bin ./scripts/install.sh --version v0.0.3

# Install with custom directory via flag
./scripts/install.sh -d /opt/tools

# Install with NO_COLOR output
NO_COLOR=1 ./scripts/install.sh
```

### Requirements

#### Dependencies
- `curl` OR `wget` (at least one required)
- `tar` (for archive extraction)
- `sha256sum` OR `shasum` (for checksum verification)

#### Supported Platforms
- **OS:** Linux, macOS
- **Architecture:** amd64 (x86_64), arm64 (aarch64)

#### Installation Prerequisites
- Write permissions to `INSTALL_DIR` (or use `sudo`)
- Network access to GitHub (github.com, api.github.com)

### Installation Process

1. Detects OS and CPU architecture
2. Validates dependencies (curl/wget, tar, checksum tools)
3. Fetches latest version from GitHub API (if not specified)
4. Downloads release archive and checksums.txt
5. Verifies SHA256 checksum of archive
6. Extracts binary from archive
7. Moves binary to `INSTALL_DIR`
8. Sets executable permissions
9. Verifies installation and shows version

### Troubleshooting

#### "Either curl or wget is required"

```bash
# macOS - Install curl via Homebrew
brew install curl

# Ubuntu/Debian - Install wget
sudo apt-get install wget

# Rocky/CentOS
sudo dnf install curl
```

#### "Cannot write to /usr/local/bin"

Use one of these solutions:

```bash
# Option 1: Use sudo
sudo ./scripts/install.sh

# Option 2: Install to user directory (recommended)
INSTALL_DIR=~/.local/bin ./scripts/install.sh

# Option 3: Create writable directory
mkdir -p ~/bin
./scripts/install.sh --dir ~/bin
export PATH="$HOME/bin:$PATH"
```

#### "Archive not found in releases"

Script failed to download archive. Check:
- Network connectivity to github.com
- GitHub API rate limiting (429 errors)
- Specified version exists: https://github.com/cesc1802/migration-tool/releases

Try specifying a known version:
```bash
./scripts/install.sh --version v0.0.3
```

#### Path Issues After Installation

If `migrate-tool` command not found:

```bash
# Check if installed
ls -l /usr/local/bin/migrate-tool

# Add to PATH if needed
export PATH="/usr/local/bin:$PATH"

# Make permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
```

### Verification

After installation:

```bash
# Check version
migrate-tool version

# Show help
migrate-tool --help
```

### Direct Binary Download

If you prefer manual installation:

1. Visit [GitHub Releases](https://github.com/cesc1802/migration-tool/releases)
2. Download binary for your platform (e.g., `migrate-tool_0.0.3_linux_amd64.tar.gz`)
3. Download `checksums.txt` from same release
4. Verify checksum:
   ```bash
   sha256sum -c checksums.txt --ignore-missing
   ```
5. Extract and install:
   ```bash
   tar -xzf migrate-tool_0.0.3_linux_amd64.tar.gz
   sudo mv migrate-tool /usr/local/bin/
   sudo chmod +x /usr/local/bin/migrate-tool
   ```

### CI/CD Integration

For automated deployments:

```bash
#!/bin/bash
# Download and install latest version
curl -fsSL https://raw.githubusercontent.com/cesc1802/migration-tool/master/scripts/install.sh | bash

# Or with specific version
curl -fsSL https://raw.githubusercontent.com/cesc1802/migration-tool/master/scripts/install.sh | bash -s -- --version v0.0.3
```

### Build from Source

Alternative to using the install script:

```bash
git clone https://github.com/cesc1802/migration-tool.git
cd migration-tool
go install ./cmd/migrate-tool
```

Requires: Go 1.25.1 or higher

---

## Distribution Strategy

### Supported Installation Methods

| Method | OS Support | Effort | Reliability |
|--------|-----------|--------|------------|
| Download Binary | Linux, macOS, Windows | None | High |
| Install Script | Linux, macOS | Low | High |
| `go install` | All | Low | High |
| Source Build | All | Medium | Very High |

### Script Maintenance

- **Location:** `./scripts/install.sh`
- **Last Updated:** Phase 01 completion
- **Tested Platforms:** Linux (x86_64, arm64), macOS (x86_64, arm64)

---

## Security Considerations

### Checksum Verification

All releases include `checksums.txt` with SHA256 hashes. The install script automatically verifies checksums before installation.

### HTTPS Only

- GitHub API and release downloads use HTTPS
- No fallback to HTTP
- Curl/wget configured with `-fsSL` (fail on HTTP error, show errors, silent, follow redirects)

### Minimal Dependencies

- No external tools beyond standard Unix utilities (curl/wget, tar, sha256sum)
- No package manager dependencies
- No requirement for sudo for user-local installations

---

## Version Management

### Fetching Latest Version

The script uses GitHub API to determine the latest release:

```bash
curl -fsSL https://api.github.com/repos/cesc1802/migration-tool/releases/latest
```

### Version Format

Versions follow semantic versioning: `vX.Y.Z` or `vX.Y.Z-suffix`

### Pinning Versions

For reproducible deployments, always specify a version:

```bash
./scripts/install.sh --version v0.0.3
```

---

## Uninstallation

```bash
# System-wide installation
sudo rm /usr/local/bin/migrate-tool

# User-local installation
rm ~/.local/bin/migrate-tool
```

---

## See Also

- [README.md](../README.md) - Quick start guide
- [Configuration Guide](./configuration.md) - Database configuration
- [Architecture Documentation](./system-architecture.md) - System design
