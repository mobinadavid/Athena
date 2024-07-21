package service

import (
	"athena/src/database/scopes"
	blockchain_explorer_service "athena/src/services/blockchain-explorer/service"
	blockchain_model "athena/src/services/blockchain/model"
	blockchain_service "athena/src/services/blockchain/service"
	"athena/src/services/payment-gateway/drivers/crypto"
	transaction_response "athena/src/services/transaction-response"
	"athena/src/services/wallet-address/model"
	"athena/src/services/wallet-address/repository"
	"athena/src/services/wallet-address/request"
	"fmt"
	"github.com/google/uuid"
)

type IWalletAddressService interface {
	GetList() (*scopes.PaginateModel, error)
	Create(request *request.CreateWalletAddressRequest) (*model.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*model.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, request *request.CreateWalletAddressRequest) (*model.WalletAddress, error)
	HandleDeposits() error
}

type WalletAddressService struct {
	IWalletAddressRepository   repository.IWalletAddressRepository
	IBlockchainService         blockchain_service.IBlockChainService
	IBlockchainExplorerService blockchain_explorer_service.IBlockchainExplorerService
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

func (service *WalletAddressService) GetByUuid(uuid *uuid.UUID) (*model.WalletAddress, error) {
	return service.IWalletAddressRepository.GetByUuid(uuid)
}

func (service *WalletAddressService) Create(request *request.CreateWalletAddressRequest) (*model.WalletAddress, error) {

	blockchain, err := service.IBlockchainService.GetByName(request.BlockchainName)
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

	blockchain, err := service.IBlockchainService.GetByName(request.BlockchainName)
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
func (service *WalletAddressService) GetTransactions(walletAddress string, blockchain *blockchain_model.Blockchain) ([]transaction_response.Response, error) {

	switch blockchain.NativeAsset {
	case "ETH":
		explorer, err := service.IBlockchainExplorerService.GetExplorerByBlockchainId(blockchain.ID)
		if err != nil {
			return nil, fmt.Errorf("error finding explorer: %w", err)
		}

		ethTransactions, err := crypto.FetchEthTransaction(walletAddress, explorer.BaseUrl, explorer.ApiKey)
		if err != nil {
			return nil, fmt.Errorf("error fetching Ethereum transactions: %w", err)
		}

		return ethTransactions, nil

	case "TRX":

	}

	return nil, nil
}

func (service *WalletAddressService) HandleDeposits() error {

	paginatedModel, err := service.GetList()
	if err != nil {
		return err
	}

	walletAddresses := paginatedModel.Items.(*[]*model.WalletAddress)

	for _, walletAddress := range *walletAddresses {

		blockchain, err := service.IBlockchainService.GetById(walletAddress.BlockchainID)
		if err != nil {

			return fmt.Errorf("failed to get blockchain for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

		txs, err := service.GetTransactions(walletAddress.WalletAddress, blockchain)
		if err != nil {

			return fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
		}
		for _, tx := range txs {
			fmt.Printf("Transaction Details:\n")
			fmt.Printf("  Nonce: %s\n", tx.Nonce)
			fmt.Printf("  BlockHash: %s\n", tx.BlockHash)
			fmt.Printf("  From: %s\n", tx.From)
			fmt.Printf("  To: %s\n", tx.To)
			fmt.Printf("  Gas: %s\n", tx.Gas)
			fmt.Printf("  GasPrice: %s\n", tx.GasPrice)
			fmt.Printf("  Confirmations: %s\n", tx.Confirmations)
			fmt.Printf("  Timestamp: %s\n", tx.Timestamp)
			fmt.Printf("  Hash: %s\n", tx.Hash)
			fmt.Printf("  BlockNumber: %s\n", tx.BlockNumber)
			fmt.Printf("  Blockchain: %s\n", tx.BlockChain)
			fmt.Println()
		}
	}
	return nil
}
