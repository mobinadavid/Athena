package blockchain_explorer

import (
	"athena/src/database"
	"athena/src/services/blockchain-explorer/controller"
	"athena/src/services/blockchain-explorer/repository"
	"athena/src/services/blockchain-explorer/service"
	blockchain_service "athena/src/services/blockchain/service"

	"github.com/google/wire"
)

var SetContainer = wire.NewSet(
	ProvideBlockchainExplorerController,
	ProvideBlockchainExplorerService,
	ProvideBlockchainExplorerRepository,
	wire.Bind(new(service.IBlockchainExplorerService), new(*service.BlockchainExplorerService)),
	wire.Bind(new(repository.IBlockchainExplorerRepository), new(*repository.BlockchainExplorerRepository)),
)

func ProvideBlockchainExplorerController(service service.IBlockchainExplorerService) *controller.BlockchainExplorerController {
	return &controller.BlockchainExplorerController{
		IBlockchainExplorerService: service,
	}
}

func ProvideBlockchainExplorerService(repository repository.IBlockchainExplorerRepository, blockchainService *blockchain_service.BlockchainService) *service.BlockchainExplorerService {
	return &service.BlockchainExplorerService{
		IBlockchainExplorerRepository: repository,
		IBlockchainService:            blockchainService,
	}
}

func ProvideBlockchainExplorerRepository(db *database.Database) *repository.BlockchainExplorerRepository {
	return &repository.BlockchainExplorerRepository{
		IDatabaseHandler: db,
	}
}
