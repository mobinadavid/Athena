package blockchain

import (
	"athena/src/database"
	"athena/src/services/blockchain/controller"
	"athena/src/services/blockchain/repository"
	"athena/src/services/blockchain/service"
	"github.com/google/wire"
)

var SetContainer = wire.NewSet(
	ProvideBlockchainController,
	ProvideBlockchainService,
	ProvideBlockchainRepository,
)

func ProvideBlockchainController(service *service.BlockchainService) *controller.BlockchainController {
	return &controller.BlockchainController{
		IBlockchainService: service,
	}
}

func ProvideBlockchainService(repository *repository.BlockchainRepository) *service.BlockchainService {
	return &service.BlockchainService{
		IBlockchainRepository: repository,
	}
}

func ProvideBlockchainRepository(db *database.Database) *repository.BlockchainRepository {
	return &repository.BlockchainRepository{
		IDatabaseHandler: db,
	}
}
