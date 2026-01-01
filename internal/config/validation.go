package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate validates the config struct and returns user-friendly errors
func Validate(c *Config) error {
	v := validator.New()

	if err := v.Struct(c); err != nil {
		validationErrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return fmt.Errorf("config validation failed: %w", err)
		}
		var errs []string
		for _, e := range validationErrs {
			errs = append(errs, formatValidationError(e))
		}
		return fmt.Errorf("config validation failed:\n  %s", strings.Join(errs, "\n  "))
	}

	// Custom: check for unexpanded env vars
	for name, env := range c.Environments {
		if strings.HasPrefix(env.DatabaseURL, "${") && strings.HasSuffix(env.DatabaseURL, "}") {
			varName := env.DatabaseURL[2 : len(env.DatabaseURL)-1]
			return fmt.Errorf("environment %q: database_url references unset env var %q", name, varName)
		}
		if strings.Contains(env.DatabaseURL, "${") {
			return fmt.Errorf("environment %q: database_url contains unexpanded variable", name)
		}
	}

	return nil
}

func formatValidationError(e validator.FieldError) string {
	field := e.StructNamespace()
	// Simplify field names for better UX
	field = strings.ReplaceAll(field, "Config.", "")

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must have at least %s items", field, e.Param())
	default:
		return fmt.Sprintf("%s failed %s validation", field, e.Tag())
	}
}
