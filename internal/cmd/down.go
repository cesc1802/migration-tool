package cmd

import (
	"fmt"

	"github.com/cesc1802/migrate-tool/internal/migrator"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var downSteps int

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Long:  `Rollback the last migration or specified number of steps.`,
	RunE:  runDown,
}

func init() {
	downCmd.Flags().IntVar(&downSteps, "steps", 1, "Number of migrations to rollback")
	rootCmd.AddCommand(downCmd)
}

func runDown(cmd *cobra.Command, args []string) error {
	mg, err := migrator.New(envName)
	if err != nil {
		return err
	}
	defer mg.Close()

	// Get status before rollback
	status, err := mg.Status()
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	if status.Applied == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Confirmation will be handled in Phase 7

	if err := mg.Down(downSteps); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Get status after rollback
	newStatus, err := mg.Status()
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	rolledBack := status.Applied - newStatus.Applied
	fmt.Printf("Rolled back %d migration(s)\n", rolledBack)
	if newStatus.Version > 0 {
		fmt.Printf("Current version: %d\n", newStatus.Version)
	} else {
		fmt.Println("Current version: none (clean slate)")
	}

	return nil
}
