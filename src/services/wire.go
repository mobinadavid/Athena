//go:build wireinject
// +build wireinject

package services

import (
	"athena/src/database"
	"athena/src/services/blockchain"
	blockchain_explorer "athena/src/services/blockchain-explorer"
	blockchain_explorer_controller "athena/src/services/blockchain-explorer/controller"
	blockchain_controller "athena/src/services/blockchain/controller"
	wallet_address "athena/src/services/wallet-address"

	"athena/src/services/wallet-address/controller"
	"github.com/google/wire"
)

type Container struct {
	BlockchainController         *blockchain_controller.BlockchainController
	WalletAddressController      *controller.WalletAddressController
	BlockchainExplorerController *blockchain_explorer_controller.BlockchainExplorerController
}

func GetContainer() *Container {
	wire.Build(
		database.GetInstance,
		blockchain.SetContainer,
		wallet_address.SetContainer,
		blockchain_explorer.SetContainer,
		wire.Struct(new(Container), "*"),
	)

	return nil
}
