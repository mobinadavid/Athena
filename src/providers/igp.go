package providers

import (
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
	"github.com/google/wire"
)

var IgpContainer = wire.NewSet(
	ProvideIGPRepository,
	ProvideIGPService,
	wire.Bind(new(services.IIGPService), new(*services.IGPService)),
	wire.Bind(new(repositories.IIGPRepository), new(*repositories.IGPRepository)),
)

func ProvideIGPRepository(db *database.Database) *repositories.IGPRepository {
	return &repositories.IGPRepository{
		IDatabaseHandler: db,
	}
}

func ProvideIGPService(igpRepository *repositories.IGPRepository) *services.IGPService {
	return &services.IGPService{
		IIGPRepository: igpRepository,
	}
}
