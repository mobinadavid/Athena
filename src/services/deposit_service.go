package services

import (
	"athena/src/models"
	"athena/src/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type IDepositService interface {
	TransactionsExist(txHash string) (bool, error)
	AddTransactionToDeposits(transaction models.Transaction) error
	HandleDeposits() error
}

type DepositService struct {
	IDepositRepository    repositories.IDepositRepository
	IWalletAddressService IWalletAddressService
}

func (service *DepositService) TransactionsExist(txHash string) (bool, error) {
	return service.IDepositRepository.TransactionsExist(txHash)

}

func (service *DepositService) AddTransactionToDeposits(transaction *models.Transaction) error {
	return service.IDepositRepository.AddTransactionToDeposits(transaction)

}

func (service *DepositService) HandleDeposits() error {
	paginatedModel, err := service.IWalletAddressService.GetAllocatedList()
	if err != nil {
		return err
	}

	walletAddresses := paginatedModel.Items.(*[]*models.WalletAddress)
	for _, walletAddress := range *walletAddresses {
		totalItems, txs, err := service.IWalletAddressService.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

		fmt.Printf("transactions found : %d\n", totalItems)
		filteredTxs := filterTransactions(txs, walletAddress)
		for _, tx := range filteredTxs {
			exists, err := service.IDepositRepository.TransactionsExist(tx.Hash)
			if err != nil {
				return fmt.Errorf("failed to check transaction existence: %w", err)
			}

			if !exists {
				if err := sendToWebhook(tx, walletAddress.WebhookURL); err != nil {
					fmt.Printf("failed to send transaction to webhook: %v\n", err)
				}

				// Add transaction to the deposits table
				err = service.IDepositRepository.AddTransactionToDeposits(tx)
				if err != nil {
					fmt.Printf("failed to add transaction to deposits: %v\n", err)
				}
			}
		}

	}
	return nil
}

func sendToWebhook(tx interface{}, webhookUrl string) error {
	// Marshal the transaction data to JSON
	txData, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	// Create a POST request
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(txData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Retry logic
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if err != nil {
			fmt.Printf("failed to send HTTP request, attempt %d: %v\n", i+1, err)
		} else {
			fmt.Printf("webhook returned status %d, attempt %d\n", resp.StatusCode, i+1)
		}
		time.Sleep(2 * time.Second) // wait before retrying
	}
	if err != nil {
		return fmt.Errorf("failed to send HTTP request after retries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}

func filterTransactions(txs []*models.Transaction, walletAddress *models.WalletAddress) []*models.Transaction {
	filteredTxs := []*models.Transaction{}
	for _, tx := range txs {
		switch walletAddress.Blockchain.NativeAsset {
		case "ETH":
			// Filter criteria for Ethereum
			if tx.IsConfirmed && tx.From == walletAddress.WalletAddress {
				amount, err := strconv.ParseFloat(tx.Amount, 64)
				if err == nil {
					tx.Amount = fmt.Sprintf("%f", amount*1e-18)
					filteredTxs = append(filteredTxs, tx)
				}
			}
		case "TRX":
			// Filter criteria for Tron
			if tx.IsConfirmed && tx.From == walletAddress.WalletAddress {
				amount, err := strconv.ParseFloat(tx.Amount, 64)
				if err == nil {
					tx.Amount = fmt.Sprintf("%f", amount*1e-6)
					filteredTxs = append(filteredTxs, tx)
				}
			}
		// Add cases for other blockchains as needed
		default:
			fmt.Printf("Unsupported blockchain: %s\n", walletAddress.Blockchain.Name)
		}
	}
	return filteredTxs
}
