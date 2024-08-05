//go:build wireinject
// +build wireinject

package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/services"
	"github.com/google/wire"
)

type Container struct {
	BlockchainController         *controllers.BlockchainController
	WalletAddressController      *controllers.WalletAddressController
	BlockchainExplorerController *controllers.BlockchainExplorerController
	DepositService               *services.DepositService
}

func GetContainer() *Container {
	wire.Build(
		database.GetInstance,
		BlockchainContainer,
		ExplorerContainer,
		WalletAddressContainer,
		DepositContainer,
		wire.Struct(new(Container), "*"),
	)

	return nil
}
