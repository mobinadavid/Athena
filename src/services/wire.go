//go:build wireinject
// +build wireinject

package services

import (
	"athena/src/database"
	"athena/src/services/blockchain"
	"athena/src/services/blockchain/controller"
	"github.com/google/wire"
)

type Container struct {
	BlockchainController *controller.BlockchainController
}

func GetContainer() *Container {
	wire.Build(
		database.GetInstance,
		blockchain.SetContainer,
		wire.Struct(new(Container), "*"),
	)

	return nil
}
