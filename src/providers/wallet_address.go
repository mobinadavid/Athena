package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
	"github.com/google/wire"
)

var WalletAddressContainer = wire.NewSet(
	ProvideWalletAddressController,
	ProvideWalletAddressService,
	ProvideWalletAddressRepository,
	wire.Bind(new(services.IWalletAddressService), new(*services.WalletAddressService)),
	wire.Bind(new(repositories.IWalletAddressRepository), new(*repositories.WalletAddressRepository)),
)

func ProvideWalletAddressController(service services.IWalletAddressService) *controllers.WalletAddressController {
	return &controllers.WalletAddressController{
		IWalletAddressService: service,
	}
}

func ProvideWalletAddressService(repository repositories.IWalletAddressRepository, blockchainService services.IBlockChainService, explorerService services.IBlockchainExplorerService) *services.WalletAddressService {
	return &services.WalletAddressService{
		IWalletAddressRepository:   repository,
		IBlockchainService:         blockchainService,
		IBlockchainExplorerService: explorerService,
	}
}

func ProvideWalletAddressRepository(db *database.Database) *repositories.WalletAddressRepository {
	return &repositories.WalletAddressRepository{
		IDatabaseHandler: db,
	}
}
