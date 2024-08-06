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
	"strconv"
	"strings"
)

type IWalletAddressService interface {
	GetList(params *scopes.QueryBuilderModel) (*scopes.PaginateModel, error)
	Create(request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *requests.CreateWalletAddressRequest) (*models.WalletAddress, error)
	GetAllocatedList() (*scopes.PaginateModel, error)
	GetTransactions(address *models.WalletAddress, page, limit uint) (*scopes.PaginateModel, error)
	GetTransactionsList(walletAddress string, blockchain *models.Blockchain) ([]*models.Transaction, error)
	AllocateWalletAddresses(request *requests.AllocateWalletAddress) ([]string, error)
	FilterTransactions(txs []*models.Transaction, walletAddress *models.WalletAddress) ([]*models.Response, error)
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
func (service *WalletAddressService) GetTransactionsList(walletAddress string, blockchain *models.Blockchain) ([]*models.Transaction, error) {
	explorerBaseUrl, err := service.IBlockchainExplorerService.GetExplorerByBlockchain(blockchain)
	if err != nil {
		return nil, fmt.Errorf("error finding explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	explorerFactory := &crypto.ExplorerFactory{}
	explorer, err := explorerFactory.CreateExplorer(blockchain, explorerBaseUrl.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating explorer for blockchain %s: %w", blockchain.NativeAsset, err)
	}

	transactions, err := explorer.FetchTransactions(walletAddress)
	if err == nil {
		return transactions, nil
	}

	return nil, fmt.Errorf("error fetching transactions %w", err)
}

func (service *WalletAddressService) GetTransactions(walletAddress *models.WalletAddress, page, limit uint) (*scopes.PaginateModel, error) {
	txs, err := service.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress, err)
	}

	filteredTxs, err := service.FilterTransactions(txs, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to filter transactions for wallet address %s: %w", walletAddress, err)
	}

	totalItems := int64(len(filteredTxs))
	totalPages := int64(math.Ceil(float64(totalItems) / float64(limit)))
	return &scopes.PaginateModel{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Items:       &filteredTxs,
	}, nil
}

func (service *WalletAddressService) FilterTransactions(txs []*models.Transaction, walletAddress *models.WalletAddress) ([]*models.Response, error) {
	filteredTxs := []*models.Response{}

	switch walletAddress.Blockchain.NativeAsset {
	case "ETH", "BSC":
		for _, tx := range txs {
			//convert confirmations to int
			confirmations, err := strconv.Atoi(tx.Confirmations)
			if err != nil {
				return nil, fmt.Errorf("error parsing Confirmations: %w", err)
			}

			if confirmations > 12 && tx.From == strings.ToLower(walletAddress.WalletAddress) {
				amount, err := strconv.ParseFloat(tx.Amount, 64)
				if err == nil {
					tx.Amount = fmt.Sprintf("%f", amount*1e-18)
				}

				// Calculate fee based on gasUsed and gasPrice
				gasUsed, err := strconv.ParseFloat(tx.GasUsed, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing GasUsed: %w", err)
				}

				gasPrice, err := strconv.ParseFloat(tx.GasPrice, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing GasPrice: %w", err)
				}

				fee := gasUsed * gasPrice * 1e-18
				tx.Fee = fmt.Sprintf("%f", fee)

				filteredTxs = append(filteredTxs, &models.Response{
					BlockNumber: tx.BlockNumber,
					Hash:        tx.Hash,
					Timestamp:   tx.Timestamp,
					From:        tx.From,
					ToAddresses: tx.ToAddresses,
					Fee:         tx.Fee,
					IsConfirmed: true,
					Value:       tx.Amount,
					BlockChain:  tx.BlockChain,
				})
			}
		}
	case "TRX":
		for _, tx := range txs {
			if tx.IsConfirmed && tx.From == walletAddress.WalletAddress {
				amount, err := strconv.ParseFloat(tx.Amount, 64)
				if err == nil {
					tx.Amount = fmt.Sprintf("%f", amount*1e-6)
				}

				fee, err := strconv.ParseFloat(tx.Fee, 64)
				if err == nil {
					tx.Fee = fmt.Sprintf("%f", fee*1e-6)
				}

				filteredTxs = append(filteredTxs, &models.Response{
					BlockNumber: tx.BlockNumber,
					Hash:        tx.Hash,
					Timestamp:   tx.Timestamp,
					From:        tx.From,
					ToAddresses: tx.ToAddresses,
					Fee:         tx.Fee,
					IsConfirmed: tx.IsConfirmed,
					Value:       tx.Amount,
					BlockChain:  tx.BlockChain,
				})
			}
		}
	case "BTC":
		for _, tx := range txs {
			filteredTxs = append(filteredTxs, &models.Response{
				BlockNumber: tx.BlockNumber,
				Hash:        tx.Hash,
				Timestamp:   tx.Timestamp,
				From:        tx.From,
				ToAddresses: tx.ToAddresses,
				Fee:         tx.Fee,
				IsConfirmed: tx.IsConfirmed,
				Value:       tx.Amount,
				BlockChain:  tx.BlockChain,
			})
		}
	default:
		return nil, fmt.Errorf("Unsupported blockchain: %s\n", walletAddress.Blockchain.Name)

	}

	return filteredTxs, nil
}
