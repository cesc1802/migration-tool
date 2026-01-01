package config

import (
	"strings"
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	c := &Config{
		Environments: map[string]Environment{
			"dev": {
				DatabaseURL:    "postgres://localhost:5432/dev",
				MigrationsPath: "./migrations",
			},
		},
	}

	err := Validate(c)
	if err != nil {
		t.Errorf("Validate() returned error for valid config: %v", err)
	}
}

func TestValidate_MissingEnvironments(t *testing.T) {
	c := &Config{
		Environments: map[string]Environment{},
	}

	err := Validate(c)
	if err == nil {
		t.Error("Validate() should return error for empty environments")
	}
	if !strings.Contains(err.Error(), "Environments") {
		t.Errorf("error should mention Environments, got: %v", err)
	}
}

func TestValidate_MissingDatabaseURL(t *testing.T) {
	c := &Config{
		Environments: map[string]Environment{
			"dev": {
				MigrationsPath: "./migrations",
			},
		},
	}

	err := Validate(c)
	if err == nil {
		t.Error("Validate() should return error for missing database_url")
	}
	if !strings.Contains(err.Error(), "DatabaseURL") {
		t.Errorf("error should mention DatabaseURL, got: %v", err)
	}
}

func TestValidate_UnexpandedEnvVar(t *testing.T) {
	c := &Config{
		Environments: map[string]Environment{
			"prod": {
				DatabaseURL:    "${UNSET_DATABASE_URL}",
				MigrationsPath: "./migrations",
			},
		},
	}

	err := Validate(c)
	if err == nil {
		t.Error("Validate() should return error for unexpanded env var")
	}
	if !strings.Contains(err.Error(), "UNSET_DATABASE_URL") {
		t.Errorf("error should mention the env var name, got: %v", err)
	}
}

func TestValidate_PartiallyUnexpandedEnvVar(t *testing.T) {
	c := &Config{
		Environments: map[string]Environment{
			"prod": {
				DatabaseURL:    "postgres://${DB_USER}@localhost:5432/db",
				MigrationsPath: "./migrations",
			},
		},
	}

	err := Validate(c)
	if err == nil {
		t.Error("Validate() should return error for partially unexpanded env var")
	}
	if !strings.Contains(err.Error(), "unexpanded") {
		t.Errorf("error should mention unexpanded variable, got: %v", err)
	}
}

func TestValidate_MultipleEnvironments(t *testing.T) {
	c := &Config{
		Environments: map[string]Environment{
			"dev": {
				DatabaseURL:    "postgres://localhost:5432/dev",
				MigrationsPath: "./migrations",
			},
			"staging": {
				DatabaseURL:         "postgres://staging:5432/staging",
				MigrationsPath:      "./migrations",
				RequireConfirmation: true,
			},
			"prod": {
				DatabaseURL:         "postgres://prod:5432/prod",
				MigrationsPath:      "./migrations",
				RequireConfirmation: true,
			},
		},
	}

	err := Validate(c)
	if err != nil {
		t.Errorf("Validate() returned error for valid multi-env config: %v", err)
	}
}
