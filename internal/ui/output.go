package ui

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorBold   = "\033[1m"
)

// UseColor checks if we should use colored output.
// Returns true if stdout is TTY and NO_COLOR env is not set.
func UseColor() bool {
	return IsTTY() && os.Getenv("NO_COLOR") == ""
}

// Success prints a success message with green checkmark
func Success(msg string) {
	if UseColor() {
		fmt.Printf("%s%s%s %s%s\n", ColorGreen, ColorBold, "OK", msg, ColorReset)
	} else {
		fmt.Printf("OK %s\n", msg)
	}
}

// Warning prints a warning message in yellow
func Warning(msg string) {
	if UseColor() {
		fmt.Printf("%s! %s%s\n", ColorYellow, msg, ColorReset)
	} else {
		fmt.Printf("! %s\n", msg)
	}
}

// Error prints an error message in red to stderr
func Error(msg string) {
	if UseColor() {
		fmt.Fprintf(os.Stderr, "%sERROR: %s%s\n", ColorRed, msg, ColorReset)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
	}
}

// Info prints an informational message in blue
func Info(msg string) {
	if UseColor() {
		fmt.Printf("%s* %s%s\n", ColorBlue, msg, ColorReset)
	} else {
		fmt.Printf("* %s\n", msg)
	}
}
