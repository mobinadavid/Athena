package service

import (
	"athena/src/database/scopes"
	blockchain_repository "athena/src/services/blockchain/repository"
	"athena/src/services/wallet-address/model"
	"athena/src/services/wallet-address/repository"
	"athena/src/services/wallet-address/request"
	"fmt"
	"github.com/google/uuid"
	"math"
)

type IWalletAddressService interface {
	GetList(page uint, limit uint) (*scopes.PaginateModel, error)
	Create(request *request.CreateWalletAddressRequest) (*model.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*model.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *request.CreateWalletAddressRequest) (*model.WalletAddress, error)
}

type WalletAddressService struct {
	IWalletAddressRepository repository.IWalletAddressRepository
	IBlockchainRepository    blockchain_repository.IBlockchainRepository
}

func (service *WalletAddressService) GetList(page uint, limit uint) (*scopes.PaginateModel, error) {

	walletAddresses, err := service.IWalletAddressRepository.GetList(page, limit)
	if err != nil {
		return nil, err
	}

	allWalletAddressCount, err := service.IWalletAddressRepository.GetCount()
	if err != nil {
		return nil, err
	}

	totalWalletAddress := int64(math.Ceil(float64(allWalletAddressCount) / float64(limit)))
	return &scopes.PaginateModel{
		Limit:       limit,
		CurrentPage: page,
		TotalPages:  totalWalletAddress,
		TotalItems:  allWalletAddressCount,
		Items:       &walletAddresses,
	}, nil

}

func (service *WalletAddressService) GetByUuid(uuid *uuid.UUID) (*model.WalletAddress, error) {
	return service.IWalletAddressRepository.GetByUuid(uuid)
}

func (service *WalletAddressService) Create(request *request.CreateWalletAddressRequest) (*model.WalletAddress, error) {
	//
	//blockchain, err := service.IBlockchainRepository.GetById(request.BlockchainId)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to find blockchain with ID %d: %w", request.BlockchainId, err)
	//}
	//// Validate wallet address format
	//if !validator.IsValidWalletAddress(blockchain.NativeAsset, request.WalletAddress) {
	//	return nil, fmt.Errorf("the walletAddress for the blockchain %s is not valid ", blockchain.NativeAsset)
	//}
	blockchain, err := service.IBlockchainRepository.GetByName(request.BlockchainName)
	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain with name %s: %w", request.BlockchainName, err)
	}

	walletAddress := &model.WalletAddress{
		WalletAddress:     request.WalletAddress,
		BlockchainID:      blockchain.ID,
		WalletAddressName: request.WalletAddressName,
		IsActive:          request.IsActive,
	}

	walletOrm, err := service.IWalletAddressRepository.Create(walletAddress)
	if err != nil {
		return nil, err
	}

	return walletOrm, nil

}

func (service *WalletAddressService) Delete(uuid *uuid.UUID) error {
	return service.IWalletAddressRepository.Delete(uuid)
}

func (service *WalletAddressService) Update(uuid *uuid.UUID, request *request.CreateWalletAddressRequest) (*model.WalletAddress, error) {

	blockchain, err := service.IBlockchainRepository.GetByName(request.BlockchainName)
	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain with name %s: %w", request.BlockchainName, err)
	}

	return service.IWalletAddressRepository.Update(uuid, &model.WalletAddress{
		WalletAddress:     request.WalletAddress,
		BlockchainID:      blockchain.ID,
		WalletAddressName: request.WalletAddressName,
		IsActive:          request.IsActive,
	})
}

// GetTransactions This method will connect to related blockchain explorer and returns the list of transactions for requested wallet address.
//func (service *WalletAddressService) GetTransactions(walletAddress string, page uint, limit uint) (*scopes.PaginateModel, error) {
//
//}
//
//func (service *WalletAddressService) HandleDeposits() error {
//	for key, walletAddress := range service.GetList() {
//		// define blockchain
//
//		txs, err := service.GetTransactions(walletAddress)
//
//	}
//}
