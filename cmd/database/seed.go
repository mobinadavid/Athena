package database

import (
	"athena/src/database/seeders"
	"log"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Manage database seeders",
}

var seedRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run all seeders",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running seeders")
		seeders.SeedBlockchain()
		seeders.SeedExplorer()
		seeders.SeedPermission()
		seeders.SeedPermissionGroup()
		seeders.SeedAuthorization()
		seeders.SeedAdmins()
		seeders.SeedUsers()
		log.Println("Database has seeded successfully!")
	},
}

func init() {
	seedCmd.AddCommand(
		seedRunCmd,
	)
}
