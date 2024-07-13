package service

import "athena/src/services/blockchain/repository"

type BlockchainService struct {
	IBlockchainRepository *repository.BlockchainRepository
}
