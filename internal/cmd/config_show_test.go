package cmd

import "testing"

func TestMaskDatabaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "postgres with password",
			input:    "postgres://user:secretpassword@localhost:5432/db",
			expected: "postgres://user:***@localhost:5432/db",
		},
		{
			name:     "mysql with password",
			input:    "mysql://admin:p@ssw0rd!@db.example.com:3306/mydb",
			expected: "mysql://admin:***@db.example.com:3306/mydb",
		},
		{
			name:     "no password",
			input:    "postgres://localhost:5432/db",
			expected: "postgres://localhost:5432/db",
		},
		{
			name:     "user without password",
			input:    "postgres://user@localhost:5432/db",
			expected: "postgres://user@localhost:5432/db",
		},
		{
			name:     "with query params",
			input:    "postgres://user:pass@localhost:5432/db?sslmode=require",
			expected: "postgres://user:***@localhost:5432/db?sslmode=require",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no scheme",
			input:    "localhost:5432/db",
			expected: "localhost:5432/db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskDatabaseURL(tt.input)
			if result != tt.expected {
				t.Errorf("maskDatabaseURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
