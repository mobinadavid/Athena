package services

import (
	"athena/src/api/http/requests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/repositories"
	"github.com/google/uuid"
	"math"
)

type IBlockChainService interface {
	GetList(page, limit uint) (*scopes.PaginateModel, error)
	Create(request *requests.CreateBlockchainRequest) (*models.Blockchain, error)
	GetByUuid(uuid *uuid.UUID) (*models.Blockchain, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateBlockchainRequest) (*models.Blockchain, error)
	GetByName(name string) (*models.Blockchain, error)
}

type BlockchainService struct {
	IBlockchainRepository repositories.IBlockchainRepository
}

func (service *BlockchainService) GetList(page, limit uint) (*scopes.PaginateModel, error) {
	blockchains, err := service.IBlockchainRepository.GetList(page, limit)
	if err != nil {
		return nil, err
	}

	allBlockchainCount, err := service.IBlockchainRepository.GetCount()
	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(allBlockchainCount) / float64(limit)))

	return &scopes.PaginateModel{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  totalPages,
		TotalItems:  allBlockchainCount,
		Items:       &blockchains,
	}, nil

}

func (service *BlockchainService) GetByUuid(uuid *uuid.UUID) (*models.Blockchain, error) {
	return service.IBlockchainRepository.GetByUuid(uuid)
}

func (service *BlockchainService) GetByName(name string) (*models.Blockchain, error) {
	return service.IBlockchainRepository.GetByName(name)
}

func (service *BlockchainService) Create(request *requests.CreateBlockchainRequest) (*models.Blockchain, error) {
	blockchain := &models.Blockchain{
		NativeAsset: request.NativeAsset,
		Title:       request.Title,
		IsActive:    request.IsActive,
		Name:        request.Name,
	}
	blockchainOrm, err := service.IBlockchainRepository.Create(blockchain)
	if err != nil {
		return nil, err
	}
	return blockchainOrm, nil
}

func (service *BlockchainService) Delete(uuid *uuid.UUID) error {
	return service.IBlockchainRepository.Delete(uuid)
}

func (service *BlockchainService) Update(uuid *uuid.UUID, request *requests.CreateBlockchainRequest) (*models.Blockchain, error) {
	blockchain := &models.Blockchain{
		NativeAsset: request.NativeAsset,
		Title:       request.Title,
		IsActive:    request.IsActive,
		Name:        request.Name,
	}
	blockchainOrm, err := service.IBlockchainRepository.Update(uuid, blockchain)
	if err != nil {
		return nil, err
	}
	return blockchainOrm, nil
}
