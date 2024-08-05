package repositories

import (
	"athena/src/database"
	"athena/src/models"
	"fmt"
)

type IDepositRepository interface {
	TransactionsExist(txHash string) (bool, error)
	AddTransactionToDeposits(transaction *models.Transaction) error
}

type DepositRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *DepositRepository) TransactionsExist(txHash string) (bool, error) {
	var count int64
	result := repository.IDatabaseHandler.GetClient().
		Model(&models.Deposits{}).
		Where("hash = ?", txHash).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check if transaction exists: %s", result.Error.Error())
	}

	return count > 0, nil
}

func (repository *DepositRepository) AddTransactionToDeposits(transaction *models.Transaction) error {
	newTransaction := &models.Deposits{
		Hash: transaction.Hash,
	}

	// Insert the new transaction into the deposits table
	result := repository.IDatabaseHandler.GetClient().Create(newTransaction)
	if result.Error != nil {
		return fmt.Errorf("failed to insert transaction into deposits table: %s", result.Error.Error())
	}
	return nil
}
