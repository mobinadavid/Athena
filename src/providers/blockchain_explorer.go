package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	blockchain_service "athena/src/services"
	"github.com/google/wire"
)

var ExplorerContainer = wire.NewSet(
	ProvideBlockchainExplorerController,
	ProvideBlockchainExplorerService,
	ProvideBlockchainExplorerRepository,
	wire.Bind(new(blockchain_service.IBlockchainExplorerService), new(*blockchain_service.BlockchainExplorerService)),
	wire.Bind(new(repositories.IBlockchainExplorerRepository), new(*repositories.BlockchainExplorerRepository)),
)

func ProvideBlockchainExplorerController(service blockchain_service.IBlockchainExplorerService) *controllers.BlockchainExplorerController {
	return &controllers.BlockchainExplorerController{
		IBlockchainExplorerService: service,
	}
}

func ProvideBlockchainExplorerService(repository repositories.IBlockchainExplorerRepository, blockchainService *blockchain_service.BlockchainService) *blockchain_service.BlockchainExplorerService {
	return &blockchain_service.BlockchainExplorerService{
		IBlockchainExplorerRepository: repository,
		IBlockchainService:            blockchainService,
	}
}

func ProvideBlockchainExplorerRepository(db *database.Database) *repositories.BlockchainExplorerRepository {
	return &repositories.BlockchainExplorerRepository{
		IDatabaseHandler: db,
	}
}
