package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// Config represents the root configuration
type Config struct {
	Environments map[string]Environment `mapstructure:"environments" validate:"required,min=1,dive"`
	Defaults     Defaults               `mapstructure:"defaults"`
}

// Environment represents per-environment configuration
type Environment struct {
	DatabaseURL         string `mapstructure:"database_url" validate:"required"`
	MigrationsPath      string `mapstructure:"migrations_path"`
	RequireConfirmation bool   `mapstructure:"require_confirmation"`
}

// Defaults represents default configuration values
type Defaults struct {
	MigrationsPath      string `mapstructure:"migrations_path"`
	RequireConfirmation bool   `mapstructure:"require_confirmation"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgErr  error
)

// Load reads config from viper, expands env vars, applies defaults, validates.
// Thread-safe: uses sync.Once to ensure single initialization.
func Load() (*Config, error) {
	cfgOnce.Do(func() {
		var c Config
		if err := viper.Unmarshal(&c); err != nil {
			cfgErr = fmt.Errorf("unmarshal config: %w", err)
			return
		}

		// Apply defaults + expand env vars
		for name, env := range c.Environments {
			if env.MigrationsPath == "" {
				env.MigrationsPath = c.Defaults.MigrationsPath
			}
			if env.MigrationsPath == "" {
				env.MigrationsPath = "./migrations"
			}
			env.DatabaseURL = ExpandEnvVars(env.DatabaseURL)
			c.Environments[name] = env
		}

		if err := Validate(&c); err != nil {
			cfgErr = err
			return
		}

		cfg = &c
	})

	return cfg, cfgErr
}

// Get returns the loaded config (nil if not loaded)
func Get() *Config {
	return cfg
}

// GetEnv retrieves environment config by name.
// Returns a copy of the environment config (modifications won't affect original).
func GetEnv(name string) (Environment, error) {
	if cfg == nil {
		return Environment{}, fmt.Errorf("config not loaded")
	}
	env, ok := cfg.Environments[name]
	if !ok {
		available := make([]string, 0, len(cfg.Environments))
		for k := range cfg.Environments {
			available = append(available, k)
		}
		return Environment{}, fmt.Errorf("environment %q not found (available: %v)", name, available)
	}
	return env, nil
}

// ResetForTesting resets the config state for testing purposes only.
// Do not use in production code.
func ResetForTesting() {
	cfg = nil
	cfgOnce = sync.Once{}
	cfgErr = nil
}
