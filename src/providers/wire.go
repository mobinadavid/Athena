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
	IpgController                *controllers.IpgController
	DepositService               *services.DepositService
	IgpService                   *services.IGPService
}

func GetContainer() *Container {
	wire.Build(
		database.GetInstance,
		BlockchainContainer,
		ExplorerContainer,
		WalletAddressContainer,
		DepositContainer,
		IpgContainer,
		IgpContainer,
		wire.Struct(new(Container), "*"),
	)

	return nil
}
