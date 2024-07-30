package services

import (
	"athena/src/api/http/requests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/repositories"
	"github.com/google/uuid"
	"math"
)

type IBlockchainExplorerService interface {
	GetList(page, limit uint) (*scopes.PaginateModel, error)
	Create(request *requests.CreateBlockchainExplorerRequest) (*models.BlockchainExplorer, error)
	GetByUuid(uuid *uuid.UUID) (*models.BlockchainExplorer, error)
	GetExplorerByBlockchain(blockchain *models.Blockchain) (*models.BlockchainExplorer, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateBlockchainExplorerRequest) (*models.BlockchainExplorer, error)
}

type BlockchainExplorerService struct {
	IBlockchainExplorerRepository repositories.IBlockchainExplorerRepository
	IBlockchainService            *BlockchainService
}

func (service *BlockchainExplorerService) GetList(page, limit uint) (*scopes.PaginateModel, error) {
	var blockchainExplorers []*models.BlockchainExplorer
	var err error
	var count int64

	blockchainExplorers, err = service.IBlockchainExplorerRepository.GetList(page, limit)
	count, err = service.IBlockchainExplorerRepository.GetCount()
	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(count) / float64(limit)))

	return &scopes.PaginateModel{
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		TotalItems:  count,
		Items:       &blockchainExplorers,
	}, nil
}

func (service *BlockchainExplorerService) GetByUuid(uuid *uuid.UUID) (*models.BlockchainExplorer, error) {
	return service.IBlockchainExplorerRepository.GetByUuid(uuid)
}

func (service *BlockchainExplorerService) GetExplorerByBlockchain(blockchain *models.Blockchain) (*models.BlockchainExplorer, error) {
	return service.IBlockchainExplorerRepository.GetExplorerByBlockchain(blockchain)
}

func (service *BlockchainExplorerService) Create(request *requests.CreateBlockchainExplorerRequest) (*models.BlockchainExplorer, error) {
	blockchainExplorer := &models.BlockchainExplorer{
		BaseUrl:   request.BaseUrl,
		Name:      request.Name,
		IsActive:  request.IsActive,
		IsDefault: request.IsDefault,
	}

	blockchains := make([]*models.Blockchain, 0, len(request.Blockchains))
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
	return service.IBlockchainExplorerRepository.Delete(uuid)
}

func (service *BlockchainExplorerService) Update(uuid *uuid.UUID, request *requests.CreateBlockchainExplorerRequest) (*models.BlockchainExplorer, error) {
	blockchainExplorer := &models.BlockchainExplorer{
		BaseUrl:   request.BaseUrl,
		Name:      request.Name,
		IsActive:  request.IsActive,
		IsDefault: request.IsDefault,
	}

	blockchains := make([]*models.Blockchain, 0, len(request.Blockchains))
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
