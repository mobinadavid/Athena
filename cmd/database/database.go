package database

import (
	"github.com/spf13/cobra"
)

// DatabaseCmd represents the base command for database operations
var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Manage the database",
	Long:  `This command allows you to manage the database, including backups, migrations, and seeders.`,
}

func init() {
	DatabaseCmd.AddCommand(
		backupCmd,
		migrateCmd,
		seedCmd,
	)
}
