package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"

	"github.com/google/wire"
)

var DepositContainer = wire.NewSet(
	ProvideDepositRepository,
	ProvideDepositService,
	ProvideDepositController,
	wire.Bind(new(services.IDepositService), new(*services.DepositService)),
	wire.Bind(new(repositories.IDepositRepository), new(*repositories.DepositRepository)),
)

func ProvideDepositRepository(db *database.Database) *repositories.DepositRepository {
	return &repositories.DepositRepository{
		IDatabaseHandler: db,
	}
}

func ProvideDepositService(
	repository repositories.IDepositRepository,
	walletAddressService services.IWalletAddressService,
	paymentRequestRepository repositories.IPaymentRequestRepository,
	notificationService services.INotificationService,
) *services.DepositService {
	return &services.DepositService{
		IDepositRepository:       repository,
		IWalletAddressService:    walletAddressService,
		PaymentRequestRepository: paymentRequestRepository,
		NotificationService:      notificationService,
	}
}

func ProvideDepositController(service services.IDepositService) *controllers.DepositController {
	return &controllers.DepositController{
		DepositService: service,
	}
}
