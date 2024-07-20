package service

import (
	"athena/src/database/scopes"
	"athena/src/services/wallet-address/model"

	"athena/src/services/wallet-address/repository"
	"athena/src/services/wallet-address/request"

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

	walletAddress := &model.WalletAddress{
		WalletAddress: request.WalletAddress,
		BlockchainID:  request.BlockchainId,
		Title:         request.Title,
		IsActive:      request.IsActive,
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
	return service.IWalletAddressRepository.Update(uuid, &model.WalletAddress{
		WalletAddress: request.WalletAddress,
		BlockchainID:  request.BlockchainId,
		Title:         request.Title,
		IsActive:      request.IsActive,
	})
}
