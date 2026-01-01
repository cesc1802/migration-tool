# Codebase Summary - migrate-tool

## Project Overview

Golang migration CLI tool - Phase 1 project setup. A cross-platform database migration tool with single-file up/down support.

**Module:** `github.com/cesc1802/migrate-tool`
**Go Version:** 1.25.1
**Total Files:** 8 files (2,194 tokens)

---

## Directory Structure

```
.
├── cmd/
│   └── migrate-tool/
│       └── main.go              # CLI entry point
├── internal/
│   └── cmd/
│       ├── root.go              # Root command & config initialization
│       └── root_test.go          # Unit tests for root command
├── Makefile                      # Build & development tasks
├── go.mod                        # Module dependencies
├── go.sum                        # Go dependencies checksums
├── migrate-tool.example.yaml     # Configuration template
├── .gitignore                    # Git ignore rules
└── .repomixignore               # Repomix ignore rules
```

---

## Core Components

### 1. CLI Entry Point (cmd/migrate-tool/main.go)
- **Purpose:** Application bootstrap and command execution
- **Key Features:**
  - Version information injection (version, commit, date)
  - Delegates command execution to internal command handler
  - Proper error handling with exit codes

### 2. Root Command Handler (internal/cmd/root.go)
- **Purpose:** Root CLI command configuration and config file handling
- **Key Features:**
  - Cobra-based command framework
  - Viper configuration management (YAML config + environment variables)
  - Persistent flags: `--config` (config file path), `--env` (environment name, default: "dev")
  - Auto-config initialization on startup

### 3. Testing (internal/cmd/root_test.go)
- **Purpose:** Unit tests for root command functionality
- **Test Coverage:**
  - `TestSetVersionInfo` - Version information setting
  - `TestExecute` - Command execution without errors
  - `TestRootCmdExists` - Root command initialization

---

## Dependencies

### Core CLI Framework
- `github.com/spf13/cobra` v1.10.2 - Command framework
- `github.com/spf13/viper` v1.21.0 - Configuration management
- `github.com/spf13/pflag` v1.0.10 - Flag parsing

### Database Migration
- `github.com/golang-migrate/migrate/v4` v4.19.1 - Migration runner

### Interactive UI
- `github.com/manifoldco/promptui` v0.9.0 - Interactive prompts

### Validation & Utilities
- `github.com/go-playground/validator/v10` v10.30.1 - Input validation
- `go.yaml.in/yaml/v3` v3.0.4 - YAML parsing

### Supporting Libraries
- `golang.org/x/crypto` v0.46.0 - Cryptography utilities
- `golang.org/x/text` v0.32.0 - Text handling
- `golang.org/x/sys` v0.39.0 - System utilities

---

## Configuration

### File: migrate-tool.example.yaml
Multi-environment configuration with database URLs and migration paths.

```yaml
environments:
  dev:
    database_url: "postgres://user:pass@localhost:5432/myapp_dev?sslmode=disable"
    migrations_path: "./migrations"
  staging:
    database_url: "${DATABASE_URL}"        # Env var support
    migrations_path: "./migrations"
    require_confirmation: true             # Safety flag
  prod:
    database_url: "${DATABASE_URL}"
    migrations_path: "./migrations"
    require_confirmation: true
```

**Config Loading:**
- Default: `./migrate-tool.yaml` in current directory
- Override: `--config <path>` flag
- Environment selection: `--env <name>` (default: dev)
- Automatic env var substitution via Viper

---

## Build & Development

### Makefile Commands

```bash
make build          # Build binary to bin/migrate-tool (cross-platform)
make run ARGS="..." # Run with arguments
make test           # Run all tests with verbose output
make lint           # Run golangci-lint
make clean          # Remove bin/ directory
```

### Build Features
- CGO disabled for cross-platform compilation
- Version information injected via ldflags:
  - `-X main.version` (VERSION, default: dev)
  - `-X main.commit` (git short hash or "none")
  - `-X main.date` (UTC timestamp)

---

## Development Patterns

### Command Structure
- Single root command: `migrate-tool`
- Persistent flags shared across subcommands
- Viper handles config file + env var merging
- Config initialized before command execution

### Error Handling
- CLI exits with code 1 on errors
- Propagates errors up to main for exit handling

### Configuration Pattern
1. Check for `--config` flag
2. Fall back to `migrate-tool.yaml` in current directory
3. Load YAML configuration
4. Merge environment variables
5. Access config via Viper during command execution

---

## Completed in Phase 1
- Project structure setup (cmd/, internal/)
- Root command with config initialization
- Makefile for build automation
- Configuration file template with environment support
- Unit tests for core functionality
- Proper gitignore setup for Go projects

---

## Next Phases
- Phase 2: Configuration schema & validation
- Phase 3: Source driver implementation
- Phase 4: Core migration commands (up, down, status)
- Phase 5-6: Utility & advanced commands
- Phase 7-8: Interactive UI & release

