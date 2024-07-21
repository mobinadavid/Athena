package blockchain

import (
	"athena/src/database"
	blockchain_explorer_service "athena/src/services/blockchain-explorer/service"
	blockchain_service "athena/src/services/blockchain/service"
	"athena/src/services/wallet-address/controller"
	"athena/src/services/wallet-address/repository"
	"athena/src/services/wallet-address/service"
	"github.com/google/wire"
)

var SetContainer = wire.NewSet(
	ProvideWalletAddressController,
	ProvideWalletAddressService,
	ProvideWalletAddressRepository,
	wire.Bind(new(service.IWalletAddressService), new(*service.WalletAddressService)),
	wire.Bind(new(repository.IWalletAddressRepository), new(*repository.WalletAddressRepository)),
)

func ProvideWalletAddressController(service service.IWalletAddressService) *controller.WalletAddressController {
	return &controller.WalletAddressController{
		IWalletAddressService: service,
	}
}

func ProvideWalletAddressService(repository repository.IWalletAddressRepository, blockchainService blockchain_service.IBlockChainService, explorerService blockchain_explorer_service.IBlockchainExplorerService) *service.WalletAddressService {
	return &service.WalletAddressService{
		IWalletAddressRepository:   repository,
		IBlockchainService:         blockchainService,
		IBlockchainExplorerService: explorerService,
	}
}

func ProvideWalletAddressRepository(db *database.Database) *repository.WalletAddressRepository {
	return &repository.WalletAddressRepository{
		IDatabaseHandler: db,
	}
}
