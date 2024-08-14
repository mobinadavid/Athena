package cmd

import (
	"athena/cmd/app"
	"athena/cmd/database"
	"athena/cmd/wallet_address"
	"athena/src/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "athena",
	Short: "Athena CLI",
	Long:  "Athena CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		app.AppCmd,
		database.DatabaseCmd,
		wallet_address.WalletAddressCmd,
	)
}

func initConfig() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
}
