package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"
	"github.com/google/wire"
)

var BlockchainContainer = wire.NewSet(
	ProvideBlockchainController,
	ProvideBlockchainService,
	ProvideBlockchainRepository,
	wire.Bind(new(services.IBlockChainService), new(*services.BlockchainService)),
	wire.Bind(new(repositories.IBlockchainRepository), new(*repositories.BlockchainRepository)),
)

func ProvideBlockchainController(service services.IBlockChainService) *controllers.BlockchainController {
	return &controllers.BlockchainController{
		IBlockchainService: service,
	}
}

func ProvideBlockchainService(repository repositories.IBlockchainRepository) *services.BlockchainService {
	return &services.BlockchainService{
		IBlockchainRepository: repository,
	}
}

func ProvideBlockchainRepository(db *database.Database) *repositories.BlockchainRepository {
	return &repositories.BlockchainRepository{
		IDatabaseHandler: db,
	}
}
