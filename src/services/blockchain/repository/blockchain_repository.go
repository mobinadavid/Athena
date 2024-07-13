package repository

import "athena/src/database"

type BlockchainRepository struct {
	IDatabaseHandler *database.Database
}
