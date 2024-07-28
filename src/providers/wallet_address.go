package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	blockchain_service "athena/src/services"
	"github.com/google/wire"
)

var WalletAddressContainer = wire.NewSet(
	ProvideWalletAddressController,
	ProvideWalletAddressService,
	ProvideWalletAddressRepository,
	wire.Bind(new(blockchain_service.IWalletAddressService), new(*blockchain_service.WalletAddressService)),
	wire.Bind(new(repositories.IWalletAddressRepository), new(*repositories.WalletAddressRepository)),
)

func ProvideWalletAddressController(service blockchain_service.IWalletAddressService) *controllers.WalletAddressController {
	return &controllers.WalletAddressController{
		IWalletAddressService: service,
	}
}

func ProvideWalletAddressService(repository repositories.IWalletAddressRepository, blockchainService blockchain_service.IBlockChainService, explorerService blockchain_service.IBlockchainExplorerService) *blockchain_service.WalletAddressService {
	return &blockchain_service.WalletAddressService{
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
