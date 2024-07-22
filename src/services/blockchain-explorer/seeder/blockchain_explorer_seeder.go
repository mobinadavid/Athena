package seeder

import (
	"athena/src/config"
	"athena/src/database"
	"athena/src/pkg/vault"
	explorerModel "athena/src/services/blockchain-explorer/model"
	blockchainModel "athena/src/services/blockchain/model"
	"context"
	"log"
)

func SeedExplorer() {
	var db = database.GetInstance()
	var configs = config.GetInstance()

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/blockchain-explorer")
	if err != nil {
		log.Println(err)
	}

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
			ApiKey:                 secrets.Data["ethApiKey"].(string),
			Blockchains:            []*blockchainModel.Blockchain{&ethereumBlockchain},
		},
		{
			IsActive:               true,
			BaseUrl:                "https://api.bscscan.com",
			BlockchainExplorerName: "bscscan",
			IsDefault:              true,
			ApiKey:                 secrets.Data["bscApiKey"].(string),
			Blockchains:            []*blockchainModel.Blockchain{&binanceBlockchain},
		},
	}

	for _, explorer := range explorers {
		db.GetClient().FirstOrCreate(&explorer, explorerModel.BlockchainExplorer{BlockchainExplorerName: explorer.BlockchainExplorerName, BaseUrl: explorer.BaseUrl})
	}

	log.Println("Blockchain_explorer Seeder executed successfully.")
}
