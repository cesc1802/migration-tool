package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

// IsTTY checks if stdout is a terminal
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Confirm prompts user for yes/no confirmation.
// Returns true if confirmed, false otherwise.
// defaultNo=true means default is "n" (safer for destructive ops).
func Confirm(message string, defaultNo bool) (bool, error) {
	if !IsTTY() {
		return false, fmt.Errorf("not a TTY: use --auto-approve for non-interactive mode")
	}

	label := message + " [y/N]"
	if !defaultNo {
		label = message + " [Y/n]"
	}

	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	answer := strings.ToLower(strings.TrimSpace(result))
	if answer == "" {
		return !defaultNo, nil
	}
	return answer == "y" || answer == "yes", nil
}

// ConfirmProduction prompts for extra confirmation for production.
// Requires typing environment name to confirm.
func ConfirmProduction(envName string) (bool, error) {
	if !IsTTY() {
		return false, fmt.Errorf("not a TTY: use --auto-approve for non-interactive mode")
	}

	fmt.Println()
	Warning(fmt.Sprintf("You are about to modify PRODUCTION environment: %s", envName))
	fmt.Println()

	// First confirmation
	prompt1 := promptui.Prompt{
		Label:     "Continue? [y/N]",
		IsConfirm: true,
	}

	if _, err := prompt1.Run(); err != nil {
		return false, nil
	}

	// Second confirmation - type environment name
	prompt2 := promptui.Prompt{
		Label: fmt.Sprintf("Type '%s' to confirm", envName),
		Validate: func(input string) error {
			if input != envName {
				return fmt.Errorf("input does not match '%s'", envName)
			}
			return nil
		},
	}

	result, err := prompt2.Run()
	if err != nil {
		return false, nil
	}

	return result == envName, nil
}

// ConfirmDangerous prompts for dangerous operations (force, goto, rollback).
// Shows warning with details before confirmation.
func ConfirmDangerous(operation, details string) (bool, error) {
	if !IsTTY() {
		return false, fmt.Errorf("not a TTY: use --auto-approve for non-interactive mode")
	}

	fmt.Println()
	fmt.Println("WARNING: DANGEROUS OPERATION")
	fmt.Println(details)
	fmt.Println()

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Proceed with %s? [y/N]", operation),
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, nil
	}

	return true, nil
}
