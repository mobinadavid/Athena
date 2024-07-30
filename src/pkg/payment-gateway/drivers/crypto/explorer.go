package crypto

import "athena/src/models"

type BlockchainExplorer interface {
	FetchTransactions(walletAddress string) ([]models.Response, error)
}
