package providers

import (
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
	"github.com/google/wire"
)

var DepositContainer = wire.NewSet(
	ProvideDepositRepository,
	ProvideDepositService,
	wire.Bind(new(services.IDepositService), new(*services.DepositService)),
	wire.Bind(new(repositories.IDepositRepository), new(*repositories.DepositRepository)),
)

func ProvideDepositRepository(db *database.Database) *repositories.DepositRepository {
	return &repositories.DepositRepository{
		IDatabaseHandler: db,
	}
}

func ProvideDepositService(repository repositories.IDepositRepository, walletAddressService services.IWalletAddressService) *services.DepositService {
	return &services.DepositService{
		IDepositRepository:    repository,
		IWalletAddressService: walletAddressService,
	}
}
