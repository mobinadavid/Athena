package services

import (
	"athena/src/api/http/requests"
	"athena/src/config"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/payment-gateway/drivers/crypto"
	"athena/src/repositories"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"reflect"
	"strconv"
	"time"
)

type IWalletAddressService interface {
	GetList(page, limit uint) (*scopes.PaginateModel, error)
	Create(request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	HandleDeposits() error
	GetActiveList() (*scopes.PaginateModel, error)
	GetTransactions(address *models.WalletAddress, page, limit uint) (*scopes.PaginateModel, error)
	AllocateWalletAddresses(request *requests.AllocateWalletAddress) ([]string, error)
}

type WalletAddressService struct {
	IWalletAddressRepository   repositories.IWalletAddressRepository
	IBlockchainService         IBlockChainService
	IBlockchainExplorerService IBlockchainExplorerService
}

func (service *WalletAddressService) GetList(page, limit uint) (*scopes.PaginateModel, error) {
	walletAddresses, err := service.IWalletAddressRepository.GetList(page, limit)
	if err != nil {
		return nil, err
	}

	allWalletAddressCount, err := service.IWalletAddressRepository.GetCount()
	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(allWalletAddressCount) / float64(limit)))

	return &scopes.PaginateModel{
		Limit:       limit,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  allWalletAddressCount,
		Items:       &walletAddresses,
	}, nil

}

func (service *WalletAddressService) AllocateWalletAddresses(request *requests.AllocateWalletAddress) ([]string, error) {
	if request.Count <= 0 {
		return nil, errors.New("invalid number of wallet address requested")
	}

	walletAddresses, err := service.IWalletAddressRepository.GetUnallocatedWalletAddress(request.Count)
	if err != nil {
		return nil, err
	}

	// Filter wallet addresses based on the blockchain name
	var filteredWalletAddresses []*models.WalletAddress
	for _, walletAddress := range walletAddresses {
		if walletAddress.Blockchain.Name == request.Blockchain {
			filteredWalletAddresses = append(filteredWalletAddresses, walletAddress)
		}
	}

	// Check if any addresses are found
	if len(filteredWalletAddresses) == 0 {
		return nil, errors.New("no wallet addresses found")
	} else {
		//Update to allocated
		err := service.IWalletAddressRepository.UpdateWalletAddressToAllocated(filteredWalletAddresses)
		if err != nil {
			return nil, err
		}
	}

	// Collect wallet address names
	var walletAddressesName []string
	for _, walletAddress := range filteredWalletAddresses {
		walletAddressesName = append(walletAddressesName, walletAddress.WalletAddress)
	}

	return walletAddressesName, nil
}

func (service *WalletAddressService) GetActiveList() (*scopes.PaginateModel, error) {
	walletAddresses, err := service.IWalletAddressRepository.GetActiveList()
	if err != nil {
		return nil, err
	}

	allWalletAddressCount, err := service.IWalletAddressRepository.GetCount()
	if err != nil {
		return nil, err
	}

	return &scopes.PaginateModel{
		TotalItems: allWalletAddressCount,
		Items:      &walletAddresses,
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
func (service *WalletAddressService) GetTransactionsList(walletAddress string, blockchain *models.Blockchain, page, limit uint) ([]models.Response, error) {
	explorerBaseUrl, err := service.IBlockchainExplorerService.GetExplorerByBlockchain(blockchain)
	if err != nil {
		return nil, fmt.Errorf("error finding explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	explorerFactory := &crypto.ExplorerFactory{}
	explorer, err := explorerFactory.CreateExplorer(blockchain, explorerBaseUrl.BaseUrl, page, limit)
	if err != nil {
		return nil, fmt.Errorf("error creating explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	maxRetry, _ := strconv.Atoi(config.GetInstance().Get("GET_TRANSACTIONS_MAX_RETRY"))
	for i := 0; i < maxRetry; i++ {
		transactions, err := explorer.FetchTransactions(walletAddress)
		if err == nil {
			return transactions, nil
		}
		fmt.Printf("Attempt %d: Error fetching transactions: %v\n", i+1, err)
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("error fetching transactions after %d attempts: %w", maxRetry, err)
}

func (service *WalletAddressService) GetTransactions(walletAddress *models.WalletAddress, page, limit uint) (*scopes.PaginateModel, error) {
	txs, err := service.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress, err)
	}

	//// Serialize transactions to JSON
	//txsJSON, err := json.Marshal(txs)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to marshal transactions for wallet address %v: %w", walletAddress, err)
	//}
	//
	//// Publish transactions to the queue
	//err = Publisher.NewPublisher(job.LogQueue).Publish(txsJSON)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to publish transactions for wallet address %v: %w", walletAddress, err)
	//}

	return &scopes.PaginateModel{
		CurrentPage: page,
		Limit:       limit,
		Items:       &txs,
	}, nil
}

func (service *WalletAddressService) HandleDeposits() error {
	paginatedModel, err := service.GetActiveList()
	if err != nil {
		return err
	}

	walletAddresses := paginatedModel.Items.(*[]*models.WalletAddress)

	for _, walletAddress := range *walletAddresses {
		txs, err := service.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

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
