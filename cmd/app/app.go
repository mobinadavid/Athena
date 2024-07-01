package app

import (
	"github.com/spf13/cobra"
)

// AppCmd Commands for interacting with apps
var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Commands for interacting with apps.",
}

func init() {
	AppCmd.AddCommand(
		bootstrapCmd,
	)
}
