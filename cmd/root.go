package cmd

import (
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

	rootCmd.AddCommand()
}

func initConfig() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
}
