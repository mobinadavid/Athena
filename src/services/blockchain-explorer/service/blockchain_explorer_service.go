package service

import (
	"athena/src/database/scopes"
	"athena/src/services/blockchain-explorer/model"
	"athena/src/services/blockchain-explorer/repository"
	"athena/src/services/blockchain-explorer/request"
	blockchain_model "athena/src/services/blockchain/model"
	"athena/src/services/blockchain/service"
	"github.com/google/uuid"
	"math"
)

type IBlockchainExplorerService interface {
	GetList(page uint, limit uint) (*scopes.PaginateModel, error)
	Create(request *request.CreateBlockchainExplorerRequest) (*model.BlockchainExplorer, error)
	GetByUuid(uuid *uuid.UUID) (*model.BlockchainExplorer, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *request.CreateBlockchainExplorerRequest) (*model.BlockchainExplorer, error)
}

type BlockchainExplorerService struct {
	IBlockchainExplorerRepository repository.IBlockchainExplorerRepository
	IBlockchainService            *service.BlockchainService
}

func (service *BlockchainExplorerService) GetList(page uint, limit uint) (*scopes.PaginateModel, error) {

	var blockchainExplorers []*model.BlockchainExplorer
	var err error
	var count int64

	blockchainExplorers, err = service.IBlockchainExplorerRepository.GetList(page, limit)
	count, err = service.IBlockchainExplorerRepository.GetCount()

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(count) / float64(limit)))

	return &scopes.PaginateModel{
		Limit:       limit,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  count,
		Items:       &blockchainExplorers,
	}, nil

}

func (service *BlockchainExplorerService) GetByUuid(uuid *uuid.UUID) (*model.BlockchainExplorer, error) {
	return service.IBlockchainExplorerRepository.GetByUuid(uuid)
}

func (service *BlockchainExplorerService) Create(request *request.CreateBlockchainExplorerRequest) (*model.BlockchainExplorer, error) {

	blockchainExplorer := &model.BlockchainExplorer{
		BaseUrl:                request.BaseUrl,
		BlockchainExplorerName: request.BlockchainExplorerName,
		IsActive:               request.IsActive,
		IsDefault:              request.IsDefault,
	}

	blockchains := make([]*blockchain_model.Blockchain, 0, len(request.Blockchains))
	for _, blockchainName := range request.Blockchains {
		blockchain, err := service.IBlockchainService.GetByName(blockchainName)
		if err != nil {
			return nil, err
		}
		blockchains = append(blockchains, blockchain)
	}

	blockchainExplorer.Blockchains = blockchains
	explorerOrm, err := service.IBlockchainExplorerRepository.Create(blockchainExplorer)
	if err != nil {
		return nil, err
	}

	return explorerOrm, nil
}

func (service *BlockchainExplorerService) Delete(uuid *uuid.UUID) error {
	return service.IBlockchainService.Delete(uuid)
}

func (service *BlockchainExplorerService) Update(uuid *uuid.UUID, request *request.CreateBlockchainExplorerRequest) (*model.BlockchainExplorer, error) {

	blockchainExplorer := &model.BlockchainExplorer{
		BaseUrl:                request.BaseUrl,
		BlockchainExplorerName: request.BlockchainExplorerName,
		IsActive:               request.IsActive,
		IsDefault:              request.IsDefault,
	}

	blockchains := make([]*blockchain_model.Blockchain, 0, len(request.Blockchains))
	for _, blockchainName := range request.Blockchains {
		blockchain, err := service.IBlockchainService.GetByName(blockchainName)
		if err != nil {
			return nil, err
		}
		blockchains = append(blockchains, blockchain)
	}

	blockchainExplorer.Blockchains = blockchains
	return service.IBlockchainExplorerRepository.Update(uuid, blockchainExplorer)

}
