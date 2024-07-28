package service

import (
	"athena/src/database/scopes"
	blockchainExplorerService "athena/src/services/blockchain-explorer/service"
	blockchainModel "athena/src/services/blockchain/model"
	blockchainService "athena/src/services/blockchain/service"
	"athena/src/services/payment-gateway/drivers/crypto/bscscan"
	"athena/src/services/payment-gateway/drivers/crypto/etherscan"
	"athena/src/services/payment-gateway/drivers/crypto/tronscan"
	"athena/src/services/wallet-address/model"
	"athena/src/services/wallet-address/repository"
	"athena/src/services/wallet-address/request"
	transactionResponse "athena/src/services/wallet-address/transaction-response"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

type IWalletAddressService interface {
	GetList() (*scopes.PaginateModel, error)
	Create(request *request.CreateWalletAddressRequest) (*model.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (map[string]interface{}, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *request.CreateWalletAddressRequest) (*model.WalletAddress, error)
	HandleDeposits() error
	GetActiveList() (*scopes.PaginateModel, error)
	GetTransactions(address map[string]interface{}) ([]transactionResponse.Response, error)
	GetWalletAddress(request *request.GetWalletAddress) ([]string, error)
}

type WalletAddressService struct {
	IWalletAddressRepository   repository.IWalletAddressRepository
	IBlockchainService         blockchainService.IBlockChainService
	IBlockchainExplorerService blockchainExplorerService.IBlockchainExplorerService
}

func (service *WalletAddressService) GetList() (*scopes.PaginateModel, error) {
	var results []map[string]interface{}

	walletAddresses, err := service.IWalletAddressRepository.GetList()
	if err != nil {
		return nil, err
	}

	for _, walletAddress := range walletAddresses {
		blockchain, err := service.IBlockchainService.GetById(walletAddress.BlockchainID)
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"id":             walletAddress.ID,
			"uuid":           walletAddress.UUID,
			"name":           walletAddress.Name,
			"wallet_address": walletAddress.WalletAddress,
			"webhook_url":    walletAddress.WebhookURL,
			"is_active":      walletAddress.IsActive,
			"allocated_at":   walletAddress.AllocatedAt,
			"created_at":     walletAddress.CreatedAt,
			"updated_at":     walletAddress.UpdatedAt,
			"deleted_at":     walletAddress.DeletedAt,
			"blockchain":     blockchain.Name,
		}
		results = append(results, result)
	}

	allWalletAddressCount, err := service.IWalletAddressRepository.GetCount()
	if err != nil {
		return nil, err
	}

	return &scopes.PaginateModel{
		TotalItems: allWalletAddressCount,
		Items:      results,
	}, nil

}

func (service *WalletAddressService) GetWalletAddress(request *request.GetWalletAddress) ([]string, error) {
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

func (service *WalletAddressService) GetByUuid(uuid *uuid.UUID) (map[string]interface{}, error) {
	walletAddress, err := service.IWalletAddressRepository.GetByUuid(uuid)
	if err != nil {
		return nil, err
	}
	blockchain, err := service.IBlockchainService.GetById(walletAddress.BlockchainID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":              walletAddress.ID,
		"uuid":            walletAddress.UUID,
		"name":            walletAddress.Name,
		"wallet_address":  walletAddress.WalletAddress,
		"webhook_url":     walletAddress.WebhookURL,
		"is_active":       walletAddress.IsActive,
		"allocated_at":    walletAddress.AllocatedAt,
		"created_at":      walletAddress.CreatedAt,
		"updated_at":      walletAddress.UpdatedAt,
		"deleted_at":      walletAddress.DeletedAt,
		"blockchain_name": blockchain.Name,
	}, nil
}

func (service *WalletAddressService) Create(request *request.CreateWalletAddressRequest) (*model.WalletAddress, error) {
	blockchain, err := service.IBlockchainService.GetByName(request.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain with name %s: %w", request.Blockchain, err)
	}

	walletAddress := &model.WalletAddress{
		WalletAddress: request.WalletAddress,
		BlockchainID:  blockchain.ID,
		Name:          request.Name,
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
	blockchain, err := service.IBlockchainService.GetByName(request.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to find blockchain with name %s: %w", request.Blockchain, err)
	}

	return service.IWalletAddressRepository.Update(uuid, &model.WalletAddress{
		WalletAddress: request.WalletAddress,
		BlockchainID:  blockchain.ID,
		Name:          request.Name,
		IsActive:      request.IsActive,
	})
}

// GetTransactionsList fetches transactions from the specified blockchain explorer
func (service *WalletAddressService) GetTransactionsList(walletAddress string, blockchain *blockchainModel.Blockchain) ([]transactionResponse.Response, error) {
	explorer, err := service.IBlockchainExplorerService.GetExplorerByBlockchainId(blockchain.ID)
	if err != nil {
		return nil, fmt.Errorf("error finding explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	var transactions []transactionResponse.Response

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

	walletAddresses := paginatedModel.Items.(*[]*model.WalletAddress)

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
func (service *WalletAddressService) GetTransactions(walletAddress map[string]interface{}) ([]transactionResponse.Response, error) {
	// Type assertions to extract values from the map
	name, ok := walletAddress["blockchain_name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'name' in walletAddress map")
	}

	address, ok := walletAddress["wallet_address"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'wallet_address' in walletAddress map")
	}

	// Fetch blockchain by ID
	blockchain, err := service.IBlockchainService.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain for wallet address %s: %w", address, err)
	}

	// Get transactions list
	txs, err := service.GetTransactionsList(address, blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for wallet address %s: %w", address, err)
	}

	return txs, nil
}
