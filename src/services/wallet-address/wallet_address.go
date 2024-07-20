package blockchain

import (
	"athena/src/database"
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

func ProvideWalletAddressService(repository repository.IWalletAddressRepository) *service.WalletAddressService {
	return &service.WalletAddressService{
		IWalletAddressRepository: repository,
	}
}

func ProvideWalletAddressRepository(db *database.Database) *repository.WalletAddressRepository {
	return &repository.WalletAddressRepository{
		IDatabaseHandler: db,
	}
}
