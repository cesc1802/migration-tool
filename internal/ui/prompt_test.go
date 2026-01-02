package ui

import (
	"testing"
)

func TestIsTTY(t *testing.T) {
	// In test environment, stdout is typically not a TTY
	result := IsTTY()
	// Just verify it doesn't panic and returns a boolean
	// In most test environments, this should be false
	if result {
		t.Log("Running in TTY environment")
	} else {
		t.Log("Running in non-TTY environment (expected in tests)")
	}
}

func TestConfirm_NonTTY(t *testing.T) {
	// In non-TTY environment, Confirm should return error
	result, err := Confirm("Test?", true)
	if err == nil {
		t.Error("Confirm should return error in non-TTY environment")
	}
	if result {
		t.Error("Confirm should return false in non-TTY environment")
	}
	if err != nil && err.Error() != "not a TTY: use --auto-approve for non-interactive mode" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestConfirmProduction_NonTTY(t *testing.T) {
	// In non-TTY environment, ConfirmProduction should return error
	result, err := ConfirmProduction("prod")
	if err == nil {
		t.Error("ConfirmProduction should return error in non-TTY environment")
	}
	if result {
		t.Error("ConfirmProduction should return false in non-TTY environment")
	}
}

func TestConfirmDangerous_NonTTY(t *testing.T) {
	// In non-TTY environment, ConfirmDangerous should return error
	result, err := ConfirmDangerous("force", "Force setting version")
	if err == nil {
		t.Error("ConfirmDangerous should return error in non-TTY environment")
	}
	if result {
		t.Error("ConfirmDangerous should return false in non-TTY environment")
	}
}
