package wallet_address

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
)

func ProvideWalletAddressController(service *service.WalletAddressService) *controller.WalletAddressController {
	return &controller.WalletAddressController{
		IWalletAddressService: service,
	}
}

func ProvideWalletAddressService(repository *repository.WalletAddressRepository) *service.WalletAddressService {
	return &service.WalletAddressService{
		IWalletAddressRepository: repository,
	}
}

func ProvideWalletAddressRepository(db *database.Database) *repository.WalletAddressRepository {
	return &repository.WalletAddressRepository{
		IDatabaseHandler: db,
	}
}
