package cmd

import (
	"fmt"

	"github.com/cesc1802/migrate-tool/internal/migrator"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var upSteps int

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply pending migrations",
	Long:  `Apply all pending migrations or specified number of steps.`,
	RunE:  runUp,
}

func init() {
	upCmd.Flags().IntVar(&upSteps, "steps", 0, "Number of migrations to apply (0 = all)")
	rootCmd.AddCommand(upCmd)
}

func runUp(cmd *cobra.Command, args []string) error {
	mg, err := migrator.New(envName)
	if err != nil {
		return err
	}
	defer mg.Close()

	// Get status before migration
	status, err := mg.Status()
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	if status.Pending == 0 {
		fmt.Println("No migrations to apply")
		return nil
	}

	// Confirmation will be handled in Phase 7
	// For now, proceed directly

	if err := mg.Up(upSteps); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	// Get status after migration
	newStatus, err := mg.Status()
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	applied := newStatus.Applied - status.Applied
	fmt.Printf("Applied %d migration(s) successfully\n", applied)
	fmt.Printf("Current version: %d\n", newStatus.Version)

	return nil
}
