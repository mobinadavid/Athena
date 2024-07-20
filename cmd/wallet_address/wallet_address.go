package wallet_address

import "github.com/spf13/cobra"

var WalletAddressCmd = &cobra.Command{
	Use: "wallet-address",
}

func init() {
	WalletAddressCmd.AddCommand(traceCmd)
}
