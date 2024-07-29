package services

import (
	"athena/src/api/http/requests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/payment-gateway/drivers/crypto/bscscan"
	"athena/src/pkg/payment-gateway/drivers/crypto/etherscan"
	"athena/src/pkg/payment-gateway/drivers/crypto/tronscan"
	"athena/src/repositories"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

type IWalletAddressService interface {
	GetList() (*scopes.PaginateModel, error)
	Create(request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	HandleDeposits() error
	GetActiveList() (*scopes.PaginateModel, error)
	GetTransactions(address *models.WalletAddress) ([]models.Response, error)
	GetWalletAddress(request *requests.GetWalletAddress) ([]string, error)
}

type WalletAddressService struct {
	IWalletAddressRepository   repositories.IWalletAddressRepository
	IBlockchainService         IBlockChainService
	IBlockchainExplorerService IBlockchainExplorerService
}

func (service *WalletAddressService) GetList() (*scopes.PaginateModel, error) {

	walletAddresses, err := service.IWalletAddressRepository.GetList()
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

func (service *WalletAddressService) GetWalletAddress(request *requests.GetWalletAddress) ([]string, error) {
	if request.Count <= 0 {
		return nil, errors.New("invalid number of wallet address requested")
	}

	walletAddresses, err := service.IWalletAddressRepository.GetWalletAddress(request.Blockchain, request.Count)
	if err != nil {
		return nil, err
	}

	return walletAddresses, nil
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
func (service *WalletAddressService) GetTransactionsList(walletAddress string, blockchain *models.Blockchain) ([]models.Response, error) {
	explorer, err := service.IBlockchainExplorerService.GetExplorerByBlockchain(blockchain)
	if err != nil {
		return nil, fmt.Errorf("error finding explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	var transactions []models.Response

	switch blockchain.NativeAsset {
	case "ETH":
		etherScanApi, err := etherscan.NewEtherscan(explorer.BaseUrl)
		if err != nil {
			return nil, fmt.Errorf("error initializing Etherscan API: %w", err)
		}
		transactions, err = etherScanApi.FetchEthTransaction(walletAddress)
		if err != nil {
			return nil, fmt.Errorf("error fetching Ethereum transactions: %w", err)
		}

	case "TRX":
		tronScanApi, err := tronscan.NewTronscan(explorer.BaseUrl)
		if err != nil {
			return nil, fmt.Errorf("error initializing TronScan API: %w", err)
		}
		transactions, err = tronScanApi.FetchTronTransactions(walletAddress)
		if err != nil {
			return nil, fmt.Errorf("error fetching Tron transactions: %w", err)
		}

	case "BSC":
		bscScanApi, err := bscscan.NewBscscan(explorer.BaseUrl)
		if err != nil {
			return nil, fmt.Errorf("error initializing BscScan API: %w", err)
		}
		transactions, err = bscScanApi.FetchBscTransaction(walletAddress)
		if err != nil {
			return nil, fmt.Errorf("error fetching BSC transactions: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported blockchain: %s", blockchain.NativeAsset)
	}

	return transactions, nil
}

func (service *WalletAddressService) HandleDeposits() error {
	paginatedModel, err := service.GetActiveList()
	if err != nil {
		return err
	}

	walletAddresses := paginatedModel.Items.(*[]*models.WalletAddress)

	for _, walletAddress := range *walletAddresses {

		blockchain, err := service.IBlockchainService.GetById(walletAddress.BlockchainID)
		if err != nil {

			return fmt.Errorf("failed to get blockchain for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

		txs, err := service.GetTransactionsList(walletAddress.WalletAddress, blockchain)
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

// Get Transactions by wallet-address
func (service *WalletAddressService) GetTransactions(walletAddress *models.WalletAddress) ([]models.Response, error) {
	blockchain, err := service.IBlockchainService.GetById(walletAddress.BlockchainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain for wallet address %s: %w", walletAddress, err)
	}

	// Get transactions list
	txs, err := service.GetTransactionsList(walletAddress.WalletAddress, blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress, err)
	}

	return txs, nil
}
