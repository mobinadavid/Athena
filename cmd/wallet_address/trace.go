package wallet_address

import (
	"athena/src/pkg/vault"
	"athena/src/providers"
	"github.com/spf13/cobra"
	"log"
)

var traceCmd = &cobra.Command{
	Use: "trace",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vault.Init(); err != nil {
			log.Fatalln(err)
		}
		serviceContainer := providers.GetContainer()
		// Call HandleDeposits method
		if err := serviceContainer.WalletAddressController.IWalletAddressService.HandleDeposits(); err != nil {
			log.Fatalf("HandleDeposits failed: %v", err)
		}

	},
}
