package database

import (
	"athena/src/pkg/vault"
	"athena/src/services/blockchain-explorer/seeder"
	seeder2 "athena/src/services/blockchain/seeder"
	"github.com/spf13/cobra"
	"log"
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
		if err := vault.Init(); err != nil {
			log.Fatalln(err)
		}
		seeder2.SeedBlockchain()
		seeder.SeedExplorer()
		log.Println("Database has seeded successfully!")
	},
}

func init() {
	seedCmd.AddCommand(
		seedRunCmd,
	)
}
