package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/repositories"
	"athena/src/services"

	"github.com/google/wire"
)

var DashboardContainer = wire.NewSet(
	ProvideDashboardService,
	ProvideDashboardController,
	wire.Bind(new(services.IDashboardService), new(*services.DashboardService)),
)

func ProvideDashboardService(
	depositRepository repositories.IDepositRepository,
	walletAddressRepository repositories.IWalletAddressRepository,
	notificationRepository repositories.INotificationRepository,
) *services.DashboardService {
	return &services.DashboardService{
		DepositRepository:       depositRepository,
		WalletAddressRepository: walletAddressRepository,
		NotificationRepository:  notificationRepository,
	}
}

func ProvideDashboardController(service services.IDashboardService) *controllers.DashboardController {
	return &controllers.DashboardController{
		DashboardService: service,
	}
}
