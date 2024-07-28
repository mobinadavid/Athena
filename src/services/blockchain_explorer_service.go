package services

import (
	"athena/src/api/http/requests"
	"athena/src/database/scopes"
	blockchain_model "athena/src/models"
	"athena/src/repositories"
	"github.com/google/uuid"
)

type IBlockchainExplorerService interface {
	GetList() (*scopes.PaginateModel, error)
	Create(request *requests.CreateBlockchainExplorerRequest) (*blockchain_model.BlockchainExplorer, error)
	GetByUuid(uuid *uuid.UUID) (*blockchain_model.BlockchainExplorer, error)
	GetExplorerByBlockchainId(blockchainId uint) (*blockchain_model.BlockchainExplorer, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateBlockchainExplorerRequest) (*blockchain_model.BlockchainExplorer, error)
}

type BlockchainExplorerService struct {
	IBlockchainExplorerRepository repositories.IBlockchainExplorerRepository
	IBlockchainService            *BlockchainService
}

func (service *BlockchainExplorerService) GetList() (*scopes.PaginateModel, error) {

	var blockchainExplorers []*blockchain_model.BlockchainExplorer
	var err error
	var count int64

	blockchainExplorers, err = service.IBlockchainExplorerRepository.GetList()
	count, err = service.IBlockchainExplorerRepository.GetCount()

	if err != nil {
		return nil, err
	}

	return &scopes.PaginateModel{
		TotalItems: count,
		Items:      &blockchainExplorers,
	}, nil

}

func (service *BlockchainExplorerService) GetByUuid(uuid *uuid.UUID) (*blockchain_model.BlockchainExplorer, error) {
	return service.IBlockchainExplorerRepository.GetByUuid(uuid)
}

func (service *BlockchainExplorerService) GetExplorerByBlockchainId(id uint) (*blockchain_model.BlockchainExplorer, error) {
	return service.IBlockchainExplorerRepository.GetExplorerByBlockchainID(id)
}

func (service *BlockchainExplorerService) Create(request *requests.CreateBlockchainExplorerRequest) (*blockchain_model.BlockchainExplorer, error) {

	blockchainExplorer := &blockchain_model.BlockchainExplorer{
		BaseUrl:   request.BaseUrl,
		Name:      request.Name,
		IsActive:  request.IsActive,
		IsDefault: request.IsDefault,
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
	return service.IBlockchainExplorerRepository.Delete(uuid)
}

func (service *BlockchainExplorerService) Update(uuid *uuid.UUID, request *requests.CreateBlockchainExplorerRequest) (*blockchain_model.BlockchainExplorer, error) {

	blockchainExplorer := &blockchain_model.BlockchainExplorer{
		BaseUrl:   request.BaseUrl,
		Name:      request.Name,
		IsActive:  request.IsActive,
		IsDefault: request.IsDefault,
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
