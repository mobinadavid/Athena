package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"

	"github.com/google/wire"
)

var NotificationContainer = wire.NewSet(
	ProvideNotificationRepository,
	ProvideNotificationService,
	ProvideNotificationController,
	wire.Bind(new(repositories.INotificationRepository), new(*repositories.NotificationRepository)),
	wire.Bind(new(services.INotificationService), new(*services.NotificationService)),
)

func ProvideNotificationRepository(db *database.Database) *repositories.NotificationRepository {
	return &repositories.NotificationRepository{
		DatabaseHandler: db,
	}
}

func ProvideNotificationService(repository repositories.INotificationRepository) *services.NotificationService {
	return &services.NotificationService{
		NotificationRepository: repository,
	}
}

func ProvideNotificationController(service services.INotificationService) *controllers.NotificationController {
	return &controllers.NotificationController{
		NotificationService: service,
	}
}
