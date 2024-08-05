package services

import (
	"athena/src/models"
	"athena/src/repositories"
)

type IDepositService interface {
	TransactionsExist(txHash string) (bool, error)
	AddTransactionToDeposits(transaction models.Response) error
}

type DepositService struct {
	IDepositRepository repositories.IDepositRepository
}

func (service *DepositService) TransactionsExist(txHash string) (bool, error) {
	return service.IDepositRepository.TransactionsExist(txHash)

}

func (service *DepositService) AddTransactionToDeposits(transaction models.Response) error {
	return service.IDepositRepository.AddTransactionToDeposits(transaction)

}
