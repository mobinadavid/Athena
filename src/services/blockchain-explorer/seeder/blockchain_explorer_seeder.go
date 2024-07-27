package seeder

import (
	"athena/src/database"
	explorerModel "athena/src/services/blockchain-explorer/model"
	blockchainModel "athena/src/services/blockchain/model"
	"log"
)

func SeedExplorer() {
	var db = database.GetInstance()

	// Retrieve the existing Blockchain records
	var tronBlockchain, ethereumBlockchain, binanceBlockchain blockchainModel.Blockchain
	db.GetClient().Where("name = ?", "Tron").First(&tronBlockchain)
	db.GetClient().Where("name = ?", "Ethereum").First(&ethereumBlockchain)
	db.GetClient().Where("name = ?", "Binance").First(&binanceBlockchain)

	// Seed the BlockchainExplorer records and associate them with the Blockchain records
	isActive := true

	explorers := []*explorerModel.BlockchainExplorer{
		{
			IsActive:    &isActive,
			BaseUrl:     "https://apilist.tronscanapi.com",
			Name:        "tronscan",
			IsDefault:   true,
			Blockchains: []*blockchainModel.Blockchain{&tronBlockchain},
		},
		{
			IsActive:    &isActive,
			BaseUrl:     "https://api.etherscan.io",
			Name:        "etherscan",
			IsDefault:   true,
			Blockchains: []*blockchainModel.Blockchain{&ethereumBlockchain},
		},
		{
			IsActive:    &isActive,
			BaseUrl:     "https://api.bscscan.com",
			Name:        "bscscan",
			IsDefault:   true,
			Blockchains: []*blockchainModel.Blockchain{&binanceBlockchain},
		},
	}

	for _, explorer := range explorers {
		db.GetClient().FirstOrCreate(&explorer, explorerModel.BlockchainExplorer{Name: explorer.Name, BaseUrl: explorer.BaseUrl})
	}

	log.Println("Blockchain_explorer Seeder executed successfully.")
}
