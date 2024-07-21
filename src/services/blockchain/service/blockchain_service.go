package service

import (
	"athena/src/database/scopes"
	"athena/src/services/blockchain/model"
	"athena/src/services/blockchain/repository"
	"athena/src/services/blockchain/request"

	"github.com/google/uuid"
)

type IBlockChainService interface {
	GetList() (*scopes.PaginateModel, error)
	Create(request *request.CreateBlockchainRequest) (*model.Blockchain, error)
	GetByUuid(uuid *uuid.UUID) (*model.Blockchain, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *request.CreateBlockchainRequest) (*model.Blockchain, error)
	GetById(id uint) (*model.Blockchain, error)
	GetByName(name string) (*model.Blockchain, error)
}

type BlockchainService struct {
	IBlockchainRepository repository.IBlockchainRepository
}

func (service *BlockchainService) GetById(id uint) (*model.Blockchain, error) {
	return service.IBlockchainRepository.GetById(id)
}

func (service *BlockchainService) GetList() (*scopes.PaginateModel, error) {

	blockchains, err := service.IBlockchainRepository.GetList()
	if err != nil {
		return nil, err
	}

	allBlockchainCount, err := service.IBlockchainRepository.GetCount()
	if err != nil {
		return nil, err
	}

	return &scopes.PaginateModel{

		TotalItems: allBlockchainCount,
		Items:      &blockchains,
	}, nil

}

func (service *BlockchainService) GetByUuid(uuid *uuid.UUID) (*model.Blockchain, error) {
	return service.IBlockchainRepository.GetByUuid(uuid)
}

func (service *BlockchainService) GetByName(name string) (*model.Blockchain, error) {
	return service.IBlockchainRepository.GetByName(name)
}

func (service *BlockchainService) Create(request *request.CreateBlockchainRequest) (*model.Blockchain, error) {
	category := &model.Blockchain{
		NativeAsset:    request.NativeAsset,
		Title:          request.Title,
		IsActive:       request.IsActive,
		BlockchainName: request.BlockchainName,
	}
	faqOrm, err := service.IBlockchainRepository.Create(category)
	if err != nil {
		return nil, err
	}
	return faqOrm, nil
}

func (service *BlockchainService) Delete(uuid *uuid.UUID) error {
	return service.IBlockchainRepository.Delete(uuid)
}

func (service *BlockchainService) Update(uuid *uuid.UUID, request *request.CreateBlockchainRequest) (*model.Blockchain, error) {
	return service.IBlockchainRepository.Update(uuid, &model.Blockchain{
		NativeAsset:    request.NativeAsset,
		Title:          request.Title,
		IsActive:       request.IsActive,
		BlockchainName: request.BlockchainName,
	})
}
