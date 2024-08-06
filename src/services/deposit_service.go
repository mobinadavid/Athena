package services

import (
	"athena/src/models"
	"athena/src/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IDepositService interface {
	TransactionsExist(txHash string) (bool, error)
	AddTransactionToDeposits(transaction *models.Response) error
	HandleDeposits() error
}

type DepositService struct {
	IDepositRepository    repositories.IDepositRepository
	IWalletAddressService IWalletAddressService
}

func (service *DepositService) TransactionsExist(txHash string) (bool, error) {
	return service.IDepositRepository.TransactionsExist(txHash)

}

func (service *DepositService) AddTransactionToDeposits(transaction *models.Response) error {
	return service.IDepositRepository.AddTransactionToDeposits(transaction)

}

func (service *DepositService) HandleDeposits() error {
	allocatedList, err := service.IWalletAddressService.GetAllocatedList()
	if err != nil {
		return err
	}

	walletAddresses := allocatedList.Items.(*[]*models.WalletAddress)
	for _, walletAddress := range *walletAddresses {
		transactions, err := service.IWalletAddressService.GetTransactionsList(walletAddress.WalletAddress, walletAddress.Blockchain)
		if err != nil {
			return fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

		filteredTxs, err := service.IWalletAddressService.FilterTransactions(transactions, walletAddress)
		if err != nil {
			return fmt.Errorf("failed to filter transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
		}

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
