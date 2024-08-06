package services

import (
	"athena/src/config"
	"athena/src/models"
	"athena/src/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
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

	configs := config.GetInstance()
	maxRetry, _ := strconv.Atoi(configs.Get("SENT_TO_WEBHOOK_MAX_RETRY"))

	// Create a Resty client with a timeout
	client := resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)

	// Send a POST request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewBuffer(txData)).
		Post(webhookUrl)

	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("webhook returned non-200 status: %d", resp.StatusCode())
	}

	return nil
}
