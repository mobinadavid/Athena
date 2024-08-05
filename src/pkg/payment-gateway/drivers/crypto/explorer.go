package crypto

import "athena/src/models"

type BlockchainExplorer interface {
	FetchTransactions(walletAddress string) (int64, []*models.Transaction, error)
}
