package controller

import "athena/src/services/blockchain/service"

type BlockchainController struct {
	IBlockchainService *service.BlockchainService
}
