package config

import (
	"os"
	"regexp"
)

var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ExpandEnvVars replaces ${VAR} patterns with environment variable values.
// If the env var is not set, the original pattern is preserved.
func ExpandEnvVars(s string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		varName := match[2 : len(match)-1] // strip ${ and }
		if val := os.Getenv(varName); val != "" {
			return val
		}
		return match // keep original if not found
	})
}
