//go:build wireinject
// +build wireinject

package providers

import (
	blockchain_controller "athena/src/api/http/controllers"
	"athena/src/database"
	"github.com/google/wire"
)

type Container struct {
	BlockchainController         *blockchain_controller.BlockchainController
	WalletAddressController      *blockchain_controller.WalletAddressController
	BlockchainExplorerController *blockchain_controller.BlockchainExplorerController
}

func GetContainer() *Container {
	wire.Build(
		database.GetInstance,
		BlockchainContainer,
		ExplorerContainer,
		WalletAddressContainer,
		wire.Struct(new(Container), "*"),
	)

	return nil
}
