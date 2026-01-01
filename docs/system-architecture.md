# System Architecture

## Golang Migration CLI Tool - Phase 1

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    User / CI/CD Pipeline                    │
└────────────────────────┬──────────────────────────────────────┘
                         │
                         ▼
                ┌────────────────────┐
                │  migrate-tool CLI  │
                │   (Binary Entry)   │
                └────────────┬───────┘
                             │
                             ▼
                ┌────────────────────────────────────┐
                │  Command Handler (Cobra)           │
                │  - Root command initialization     │
                │  - Flag parsing                    │
                │  - Subcommand routing              │
                └────────────┬───────────────────────┘
                             │
                ┌────────────┴────────────┐
                ▼                         ▼
        ┌───────────────┐        ┌──────────────────┐
        │ Config Loader │        │ Viper Config     │
        │               │        │ Management       │
        │ - YAML file   │        │ - Env vars       │
        │ - Env vars    │        │ - Defaults       │
        └───────┬───────┘        └──────┬───────────┘
                │                       │
                └───────────┬───────────┘
                            ▼
                  ┌──────────────────────┐
                  │ Configuration Model  │
                  │ - environments       │
                  │   - dev/staging/prod │
                  │   - database_url     │
                  │   - migrations_path  │
                  │   - require_confirm  │
                  └──────┬───────────────┘
                         │
                         ▼
            ┌────────────────────────────┐
            │ Core Functionality         │  (Phase 3+)
            │ - Database Drivers         │
            │ - Migration Engine         │
            │ - Execution Handlers       │
            └────────────────────────────┘
```

---

## Layered Architecture

### Layer 1: Presentation (CLI)
**Files:** `cmd/migrate-tool/main.go`

**Responsibilities:**
- Application bootstrap
- Version information injection
- Command execution entry point
- Error handling and exit codes

**Key Components:**
```
main()
├── SetVersionInfo(version, commit, date)
└── Execute() -> rootCmd.Execute()
```

---

### Layer 2: Command Handler
**Files:** `internal/cmd/root.go`, `internal/cmd/root_test.go`

**Responsibilities:**
- Root command definition
- Persistent flag management
- Configuration initialization
- Command routing

**Key Components:**
```
rootCmd (Cobra Command)
├── Use: "migrate-tool"
├── Short: "Database migration CLI tool"
├── Long: "Cross-platform database migration..."
├── PersistentFlags:
│   ├── --config (string): config file path
│   └── --env (string): environment name [default: dev]
└── Subcommands: (Phase 2+)
    ├── migrate up
    ├── migrate down
    ├── migrate status
    └── ...
```

**Flag Processing:**
```
Flag Input
    ↓
rootCmd.PersistentFlags()
    ↓
cobra.OnInitialize(initConfig)
    ↓
initConfig() runs before command
    ↓
Viper loads config & env vars
```

---

### Layer 3: Configuration Management
**Framework:** Viper (github.com/spf13/viper)

**Configuration Flow:**
```
1. Check --config flag
   ├── If set: viper.SetConfigFile(cfgFile)
   └── If not: Use default (migrate-tool.yaml)

2. Set config parameters
   ├── ConfigName: "migrate-tool"
   ├── ConfigType: "yaml"
   └── ConfigPath: "." (current directory)

3. Enable auto env var loading
   └── viper.AutomaticEnv()

4. Read configuration
   └── viper.ReadInConfig()

5. Access during command execution
   └── viper.GetString("environments.dev.database_url")
```

**Configuration Schema (migrate-tool.yaml):**
```yaml
environments:
  <env-name>:
    database_url: <connection-string>
    migrations_path: <directory-path>
    require_confirmation: <boolean>
```

**Environment Variable Override:**
- Pattern: `MIGRATE_TOOL_<KEY_PATH_UPPERCASE>`
- Example: `MIGRATE_TOOL_ENVIRONMENTS_DEV_DATABASE_URL`
- Case-insensitive key paths

---

### Layer 4: Core Functionality (Future Phases)
**Planned Modules:**
- `internal/config/` - Configuration validation & models (Phase 2)
- `internal/driver/` - Database driver interface & implementations (Phase 3)
- `internal/migration/` - Migration parsing & execution (Phase 4)
- `internal/util/` - Utilities & helpers

---

## Data Flow

### Startup Sequence
```
1. main.go main()
   │
   ├─ SetVersionInfo(version, commit, date)
   │  └─ Stores version data globally
   │
   └─ Execute()
      └─ rootCmd.Execute()
         │
         ├─ Trigger: cobra.OnInitialize()
         │  └─ initConfig()
         │     ├─ Check --config flag
         │     ├─ Set Viper config name/type/path
         │     ├─ Enable auto env vars
         │     └─ Load migrate-tool.yaml
         │
         └─ Execute root command or route to subcommand
            └─ Subcommand accesses config via viper.Get*()
```

### Configuration Access Pattern
```
Subcommand Execution
    ↓
cfgFile (from --config flag)
envName (from --env flag, default: "dev")
    ↓
Viper Configuration Merged
├─ YAML file values
├─ Environment variable overrides
└─ Compiled into in-memory map
    ↓
Access pattern:
environments.dev.database_url → viper.GetString(...)
environments.dev.migrations_path → viper.GetString(...)
```

---

## Component Relationships

### Dependencies
```
main.go
└─ depends on → internal/cmd/root.go
                ├─ depends on → github.com/spf13/cobra
                └─ depends on → github.com/spf13/viper

root.go
├─ depends on → cobra.Command (Cobra framework)
└─ depends on → viper (Configuration)
```

### External Dependencies (Go Modules)
```
Core CLI:
├─ spf13/cobra v1.10.2 (command framework)
├─ spf13/viper v1.21.0 (config management)
├─ spf13/pflag v1.0.10 (flag parsing)
└─ spf13/afero v1.15.0 (file system abstraction)

Database:
└─ golang-migrate/migrate/v4 v4.19.1

Utilities:
├─ manifoldco/promptui v0.9.0 (interactive UI)
├─ go-playground/validator/v10 v10.30.1 (validation)
└─ go.yaml.in/yaml/v3 v3.0.4 (YAML parsing)

System:
├─ golang.org/x/sys v0.39.0
├─ golang.org/x/text v0.32.0
└─ golang.org/x/crypto v0.46.0
```

---

## Phase 1 Design Decisions

### 1. Cobra for CLI Framework
**Rationale:** Industry standard, excellent for complex CLIs
**Benefit:** Easy subcommand management, automatic help, version flags

### 2. Viper for Configuration
**Rationale:** Supports YAML + env var override elegantly
**Benefit:** Multi-environment support, automatic env binding

### 3. Package Structure (cmd/ + internal/)
**Rationale:** Go conventions - cmd/ for binaries, internal/ for private packages
**Benefit:** Clear public/private separation, encapsulation

### 4. YAML Configuration Format
**Rationale:** Human-readable, hierarchical structure
**Benefit:** Multi-environment definitions, environment variable templating

### 5. Version Injection via LDFLAGS
**Rationale:** Single binary with build metadata
**Benefit:** No external version files, accurate version info from git

### 6. Configuration-First Design
**Rationale:** Support DevOps/CI/CD workflows with environment-specific configs
**Benefit:** Same binary, different configurations per environment

---

## Scalability & Extension Points

### For Phase 2-3: Adding Subcommands
```go
// Phase 2: Configuration validation
var validateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		// Access config via viper
	},
}
func init() {
	rootCmd.AddCommand(validateCmd)
}

// Phase 4: Migration commands
rootCmd.AddCommand(upCmd, downCmd, statusCmd)
```

### For Phase 3: Database Driver Interface
```go
// Driver interface for multiple database support
type Driver interface {
	Connect(dsn string) error
	RunMigration(migration *Migration) error
	Close() error
}

// Implementations: PostgreSQL, MySQL, SQLite
```

### For Phase 4: Migration Engine
```go
// Migration execution with rollback support
type MigrationRunner struct {
	driver Driver
	migrations []*Migration
}

func (m *MigrationRunner) Up(target string) error
func (m *MigrationRunner) Down(steps int) error
func (m *MigrationRunner) Status() (*MigrationStatus, error)
```

---

## Error Handling Architecture

### Error Propagation Path
```
Driver Error
    ↓
Runner Error
    ↓
Command Handler (returns error)
    ↓
Cobra (prints error)
    ↓
main.go (exits with code 1)
```

### Error Types
- Configuration errors: Invalid YAML, missing env vars
- Driver errors: Connection failed, invalid DSN
- Migration errors: Invalid SQL, constraint violations
- CLI errors: Invalid flags, missing arguments

---

## Security Architecture

### Secret Management
```
Sensitive Data (passwords, tokens)
    ↓
Environment Variables
    ↓
Viper (never logged)
    ↓
Driver (connection string passed securely)
    ↓
Database (over secure connection)
```

### Protection Mechanisms
- Environment variable support for secrets (not in YAML)
- No hardcoded credentials
- Connection string validation before execution
- Confirmation required for destructive operations

---

## Performance Considerations

### Phase 1 Performance Profile
- **Startup:** ~50-100ms (config loading + cobra initialization)
- **Memory:** ~10-20MB (CLI tool, minimal)
- **Binary Size:** ~5-10MB (depends on dependencies)

### Optimization Points for Future Phases
- Cache compiled migration scripts
- Connection pooling for database operations
- Parallel migration execution (Phase 5+)

---

## Testing Architecture

### Unit Test Organization
```
internal/cmd/root_test.go
├── TestSetVersionInfo
│   └── Tests version info injection
├── TestExecute
│   └── Tests command execution
└── TestRootCmdExists
    └── Tests root command initialization
```

### Future Test Layers
- **Unit Tests:** Individual functions & commands
- **Integration Tests:** Config + Command + Database
- **E2E Tests:** Full workflow from CLI to migration

---

## Deployment Model

### Binary Distribution
- Single stateless binary
- Platform-specific builds (Linux, macOS, Windows)
- No external dependencies (CGO disabled)

### Configuration Management
- Local YAML file per environment
- Environment variables for overrides
- No configuration server needed

### CI/CD Integration
```
Build Stage:
├─ make build (builds with version info)
└─ Output: bin/migrate-tool-<os>-<arch>

Deploy Stage:
├─ Copy binary to server
├─ Copy migrate-tool.yaml to server
└─ Run: ./migrate-tool migrate up --env prod
```

---

## Phase 1 Summary

**Completed:**
- CLI framework foundation (Cobra)
- Configuration system (Viper + YAML)
- Version injection pipeline
- Error handling scaffolding
- Build automation (Makefile)
- Unit test foundation

**Ready for Phase 2:**
- Configuration validation
- Schema definition
- Type-safe config accessors

