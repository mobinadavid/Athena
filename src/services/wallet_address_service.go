package services

import (
	"athena/src/api/http/requests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/payment-gateway/drivers/crypto"
	"athena/src/repositories"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IWalletAddressService interface {
	GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error)
	Create(request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetAllocatedList() (*scopes.PaginatedModel, error)
	GetTransactions(address *models.WalletAddress, page, limit uint) (*scopes.PaginatedModel, error)
	GetTransactionsList(walletAddress *models.WalletAddress) ([]*models.Transaction, error)
	AllocateWalletAddresses(request *requests.AllocateWalletAddress) ([]string, error)
	FilterTransactions(txs []*models.Transaction, walletAddress *models.WalletAddress) ([]*models.TransactionResponse, error)
}

type WalletAddressService struct {
	IWalletAddressRepository   repositories.IWalletAddressRepository
	IBlockchainService         IBlockChainService
	IBlockchainExplorerService IBlockchainExplorerService
}

func (service *WalletAddressService) GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error) {
	walletAddresses, count, err := service.IWalletAddressRepository.GetList(params)
	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(count) / float64(params.Limit)))

	return &scopes.PaginatedModel{
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
	} else if len(walletAddresses) != request.Count {
		return nil, errors.New("not enough wallet addresses found")
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

func (service *WalletAddressService) GetAllocatedList() (*scopes.PaginatedModel, error) {
	walletAddresses, err := service.IWalletAddressRepository.GetAllocatedList()
	if err != nil {
		return nil, err
	}

	return &scopes.PaginatedModel{
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
		WebhookURL:    request.WebhookURL,
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
		WebhookURL:    request.WebhookURL,
	})
}

// GetTransactionsList fetches transactions from the specified blockchain explorer
func (service *WalletAddressService) GetTransactionsList(walletAddress *models.WalletAddress) ([]*models.Transaction, error) {
	explorerBaseUrl, err := service.IBlockchainExplorerService.GetExplorerByBlockchain(walletAddress.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("error finding explorer for blockchain %s: %w", walletAddress.Blockchain.NativeAsset, err)
	}

	explorerFactory := &crypto.ExplorerFactory{}
	explorer, err := explorerFactory.CreateExplorer(walletAddress.Blockchain, explorerBaseUrl.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating explorer for blockchain %s: %w", walletAddress.Blockchain.NativeAsset, err)
	}

	transactions, err := explorer.FetchTransactions(walletAddress.WalletAddress)
	if err == nil {
		return transactions, nil
	}

	return nil, fmt.Errorf("error fetching transactions %w", err)
}

func (service *WalletAddressService) GetTransactions(walletAddress *models.WalletAddress, page, limit uint) (*scopes.PaginatedModel, error) {
	//to check walletAddress is allocated and active
	isAllocated, err := service.IWalletAddressRepository.IsAllocated(walletAddress)
	if err != nil {
		return nil, err
	}
	if !isAllocated {
		return nil, fmt.Errorf("wallet address %s is not allocated or active ", walletAddress.WalletAddress)
	}

	txs, err := service.GetTransactionsList(walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for wallet address %v: %w", walletAddress, err)
	}

	filteredTxs, err := service.FilterTransactions(txs, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to filter transactions for wallet address %v: %w", walletAddress, err)
	}

	//paginating transactions
	start := (page - 1) * limit
	end := start + limit

	// Ensure indices are within bounds
	if start > uint(len(filteredTxs)) {
		start = uint(len(filteredTxs))
	}
	if end > uint(len(filteredTxs)) {
		end = uint(len(filteredTxs))
	}

	// Paginate the transactions
	paginatedItems := filteredTxs[start:end]

	totalItems := int64(len(filteredTxs))
	totalPages := int64(math.Ceil(float64(totalItems) / float64(limit)))

	return &scopes.PaginatedModel{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Items:       &paginatedItems,
	}, nil
}

func (service *WalletAddressService) FilterTransactions(txs []*models.Transaction, walletAddress *models.WalletAddress) ([]*models.TransactionResponse, error) {
	var filteredTxs []*models.TransactionResponse
	timeLayout := "Jan-02-2006 03:04:05 PM UTC"

	for _, tx := range txs {
		transactionTime, err := time.Parse(timeLayout, tx.Timestamp)
		if err != nil {
			fmt.Printf("Error parsing timestamp: %v\n", err)
			continue
		}

		if !isValidTransaction(tx, walletAddress, transactionTime) {
			continue
		}

		if err := convertAmountAndFee(tx, walletAddress.Blockchain.NativeAsset); err != nil {
			return nil, err
		}

		amount, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing amount: %w", err)
		}

		fee, err := strconv.ParseFloat(tx.Fee, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing fee: %w", err)
		}

		blockNumber, err := strconv.ParseInt(tx.BlockNumber, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing blockNumber: %w", err)
		}

		filteredTxs = append(filteredTxs, &models.TransactionResponse{
			BlockNumber: blockNumber,
			Hash:        tx.Hash,
			Timestamp:   transactionTime,
			From:        tx.From,
			ToAddresses: tx.ToAddresses,
			Fee:         fee,
			IsConfirmed: tx.IsConfirmed,
			Value:       amount,
			BlockChain:  tx.BlockChain,
		})
	}

	return filteredTxs, nil
}

// isValidTransaction checks if the transaction is valid based on the blockchain type and wallet address.
func isValidTransaction(tx *models.Transaction, walletAddress *models.WalletAddress, transactionTime time.Time) bool {
	switch walletAddress.Blockchain.NativeAsset {
	case "ETH", "BSC":
		confirmations, err := strconv.Atoi(tx.Confirmations)
		if err != nil {
			fmt.Printf("Error parsing Confirmations: %v\n", err)
			return false
		}
		return confirmations > 12 && tx.From == strings.ToLower(walletAddress.WalletAddress) && transactionTime.After(walletAddress.AllocatedAt)
	case "TRX":
		return tx.IsConfirmed && tx.From == walletAddress.WalletAddress
	case "BTC":
		return tx.From == walletAddress.WalletAddress && transactionTime.After(walletAddress.AllocatedAt)
	default:
		fmt.Printf("Unsupported blockchain: %s\n", walletAddress.Blockchain.Name)
		return false
	}
}

// convertAmountAndFee converts amount and fee based on the blockchain type.
func convertAmountAndFee(tx *models.Transaction, blockchain string) error {
	switch blockchain {
	case "ETH", "BSC":
		//amount and fee should multiply 10^-18
		amount, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			return fmt.Errorf("error parsing amount: %w", err)
		}
		tx.Amount = fmt.Sprintf("%f", amount*1e-18)
		//calculate fee by multiplying gasUsed and gasPrice
		gasUsed, err := strconv.ParseFloat(tx.GasUsed, 64)
		if err != nil {
			return fmt.Errorf("error parsing GasUsed: %w", err)
		}
		gasPrice, err := strconv.ParseFloat(tx.GasPrice, 64)
		if err != nil {
			return fmt.Errorf("error parsing GasPrice: %w", err)
		}
		tx.IsConfirmed = true
		tx.Fee = fmt.Sprintf("%f", gasUsed*gasPrice*1e-18)
	case "TRX":
		//amount and fee should multiply 10^-6
		amount, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			return fmt.Errorf("error parsing amount: %w", err)
		}
		tx.Amount = fmt.Sprintf("%f", amount*1e-6)

		fee, err := strconv.ParseFloat(tx.Fee, 64)
		if err != nil {
			return fmt.Errorf("error parsing fee: %w", err)
		}
		tx.Fee = fmt.Sprintf("%f", fee*1e-6)
	case "BTC":
		//amount and fee should multiply 10^-8
		fee, err := strconv.ParseFloat(tx.Fee, 64)
		if err != nil {
			return fmt.Errorf("error parsing fee: %w", err)
		}
		tx.Fee = fmt.Sprintf("%f", fee*1e-8)
	default:
		return fmt.Errorf("unsupported blockchain for conversion: %s", blockchain)
	}
	return nil
}
