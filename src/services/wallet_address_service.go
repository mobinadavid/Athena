package services

import (
	"athena/src/api/http/requests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/payment-gateway/drivers/crypto"
	"athena/src/repositories"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"reflect"
)

type IWalletAddressService interface {
	GetList(params *scopes.QueryBuilderModel) (*scopes.PaginateModel, error)
	Create(request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	HandleDeposits() error
	GetAllocatedList() (*scopes.PaginateModel, error)
	GetTransactions(address *models.WalletAddress, page, limit uint) (*scopes.PaginateModel, error)
	AllocateWalletAddresses(request *requests.AllocateWalletAddress) ([]string, error)
}

type WalletAddressService struct {
	IWalletAddressRepository   repositories.IWalletAddressRepository
	IBlockchainService         IBlockChainService
	IBlockchainExplorerService IBlockchainExplorerService
}

func (service *WalletAddressService) GetList(params *scopes.QueryBuilderModel) (*scopes.PaginateModel, error) {
	walletAddresses, count, err := service.IWalletAddressRepository.GetList(params)
	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(count) / float64(params.Limit)))

	return &scopes.PaginateModel{
		Limit:       params.Limit,
		CurrentPage: params.Page,
		TotalPages:  totalPages,
		TotalItems:  count,
		Items:       &walletAddresses,
	}, nil

}

func (service *WalletAddressService) AllocateWalletAddresses(request *requests.AllocateWalletAddress) ([]string, error) {
	if request.Count <= 0 {
		return nil, errors.New("invalid number of wallet address requested")
	}

	blockchain, err := service.IBlockchainService.GetByName(request.Blockchain)
	if err != nil {
		return nil, err
	}

	walletAddresses, err := service.IWalletAddressRepository.GetUnallocatedWalletAddress(request.Count, blockchain.ID)
	if err != nil {
		return nil, err
	}

	// Check if any addresses are found
	if len(walletAddresses) == 0 {
		return nil, errors.New("no wallet addresses found")
	} else {
		//Update to allocated
		err := service.IWalletAddressRepository.UpdateWalletAddressToAllocated(walletAddresses)
		if err != nil {
			return nil, err
		}
	}

	// Collect wallet address names
	var walletAddressesName []string
	for _, walletAddress := range walletAddresses {
		walletAddressesName = append(walletAddressesName, walletAddress.WalletAddress)
	}

	return walletAddressesName, nil
}

func (service *WalletAddressService) GetAllocatedList() (*scopes.PaginateModel, error) {
	walletAddresses, err := service.IWalletAddressRepository.GetAllocatedList()
	if err != nil {
		return nil, err
	}

	return &scopes.PaginateModel{
		Items: &walletAddresses,
	}, nil

}

func (service *WalletAddressService) GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error) {
	return service.IWalletAddressRepository.GetByUuid(uuid)

}

func (service *WalletAddressService) Create(request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error) {
	blockchain, err := service.IBlockchainService.GetByName(request.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain with name %s: %w", request.Blockchain, err)
	}

	walletAddress := &models.WalletAddress{
		WalletAddress: request.WalletAddress,
		BlockchainID:  blockchain.ID,
		Name:          request.Name,
		IsActive:      request.IsActive,
		Blockchain:    blockchain,
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

func (service *WalletAddressService) Update(uuid *uuid.UUID, request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error) {
	blockchain, err := service.IBlockchainService.GetByName(request.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain with name %s: %w", request.Blockchain, err)
	}

	return service.IWalletAddressRepository.Update(uuid, &models.WalletAddress{
		WalletAddress: request.WalletAddress,
		BlockchainID:  blockchain.ID,
		Name:          request.Name,
		Blockchain:    blockchain,
		IsActive:      request.IsActive,
	})
}

// GetTransactionsList fetches transactions from the specified blockchain explorer
func (service *WalletAddressService) GetTransactionsList(walletAddress string, blockchain *models.Blockchain, page, limit uint) (int64, []models.Response, error) {
	explorerBaseUrl, err := service.IBlockchainExplorerService.GetExplorerByBlockchain(blockchain)
	if err != nil {
		return 0, nil, fmt.Errorf("error finding explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	explorerFactory := &crypto.ExplorerFactory{}
	explorer, err := explorerFactory.CreateExplorer(blockchain, explorerBaseUrl.BaseUrl, page, limit)
	if err != nil {
		return 0, nil, fmt.Errorf("error creating explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	totalItems, transactions, err := explorer.FetchTransactions(walletAddress)
	if err == nil {
		return totalItems, transactions, nil
	}

	return 0, nil, fmt.Errorf("error fetching transactions %w", err)
}

func (service *WalletAddressService) GetTransactions(walletAddress *models.WalletAddress, page, limit uint) (*scopes.PaginateModel, error) {
	totalItems, txs, err := service.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress, err)
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(limit)))
	return &scopes.PaginateModel{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Items:       &txs,
	}, nil
}

func (service *WalletAddressService) HandleDeposits() error {
	paginatedModel, err := service.GetAllocatedList()
	if err != nil {
		return err
	}

	walletAddresses := paginatedModel.Items.(*[]*models.WalletAddress)
	for _, walletAddress := range *walletAddresses {
		totalItems, txs, err := service.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

		fmt.Printf("transactions found : %d\n", totalItems)
		for _, tx := range txs {
			v := reflect.ValueOf(tx)
			t := v.Type()

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				value := v.Field(i).Interface()
				fmt.Printf("  %s: %v\n", field.Name, value)
			}
			fmt.Println()
		}
	}
	return nil
}
