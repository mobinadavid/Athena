package wallet_address

import (
	"athena/src/services"
	"github.com/spf13/cobra"
	"log"
)

var traceCmd = &cobra.Command{
	Use: "trace",
	Run: func(cmd *cobra.Command, args []string) {
		serviceContainer := services.GetContainer()

		// Call HandleDeposits method
		if err := serviceContainer.WalletAddressController.IWalletAddressService.HandleDeposits(); err != nil {
			log.Fatalf("HandleDeposits failed: %v", err)
		}
	},
}
