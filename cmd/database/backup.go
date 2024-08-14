package database

import (
	"fmt"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage database backups",
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all database backups",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all database backups...")
	},
}

var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database backup created.")
	},
}

var backupCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans database backup based on their creation time",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database clean.")
	},
}

var backupDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database clean.")
	},
}

func init() {
	backupCmd.AddCommand(
		backupListCmd,
		backupCreateCmd,
		backupCleanCmd,
		backupDeleteCmd,
	)
}
