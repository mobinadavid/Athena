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
	db.GetClient().Where("blockchain_name = ?", "Tron").First(&tronBlockchain)
	db.GetClient().Where("blockchain_name = ?", "Ethereum").First(&ethereumBlockchain)
	db.GetClient().Where("blockchain_name = ?", "Binance").First(&binanceBlockchain)

	// Seed the BlockchainExplorer records and associate them with the Blockchain records
	explorers := []*explorerModel.BlockchainExplorer{
		{
			IsActive:               true,
			BaseUrl:                "https://apilist.tronscanapi.com",
			BlockchainExplorerName: "tronscan",
			IsDefault:              true,
			ApiKey:                 "",
			Blockchains:            []*blockchainModel.Blockchain{&tronBlockchain},
		},
		{
			IsActive:               true,
			BaseUrl:                "https://api.etherscan.io",
			BlockchainExplorerName: "etherscan",
			IsDefault:              true,
			ApiKey:                 "WIVAQUE9C2EFPMV2PG2HTB9CVRRHS5HCK1",
			Blockchains:            []*blockchainModel.Blockchain{&ethereumBlockchain},
		},
		{
			IsActive:               true,
			BaseUrl:                "https://api.bscscan.com",
			BlockchainExplorerName: "bscscan",
			IsDefault:              true,
			ApiKey:                 "GHK25VFCJQF1QIKPFQJZ9KEI4452HQJDAM",
			Blockchains:            []*blockchainModel.Blockchain{&binanceBlockchain},
		},
	}

	for _, explorer := range explorers {
		db.GetClient().FirstOrCreate(&explorer, explorerModel.BlockchainExplorer{BlockchainExplorerName: explorer.BlockchainExplorerName, BaseUrl: explorer.BaseUrl})
	}

	log.Println("Blockchain_explorer Seeder executed successfully.")
}
