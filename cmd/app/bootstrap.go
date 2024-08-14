package app

import (
	"athena/src/bootstrap"
	"athena/src/pkg/vault"
	"github.com/spf13/cobra"
	"log"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstraps the application and it's related services.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := vault.Init(); err != nil {
			log.Fatalln(err)
		}
		if err := bootstrap.Init(); err != nil {
			log.Fatalf("Bootstrap Service: Failed to Initialize. %v", err)
		}
	},
}
