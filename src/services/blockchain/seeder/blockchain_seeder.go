package seeder

import (
	"athena/src/database"
	"athena/src/services/blockchain/model"
	"encoding/json"

	"log"
)

func SeedBlockchain() {
	titleTron := map[string]string{
		"en": "Tron",
		"fa": "ترون",
	}
	tron, err := json.Marshal(titleTron)
	if err != nil {

	}

	titleBsc := map[string]string{
		"en": "Binance",
		"fa": "بایننس",
	}
	bsc, err := json.Marshal(titleBsc)
	if err != nil {

	}

	titleEth := map[string]string{
		"en": "Ethereum",
		"fa": "اتریوم",
	}
	eth, err := json.Marshal(titleEth)
	if err != nil {

	}

	isActive := true
	blockchains := []*model.Blockchain{
		{
			IsActive:       &isActive,
			Title:          tron,
			NativeAsset:    "TRX",
			BlockchainName: "Tron",
		},
		{
			IsActive:       &isActive,
			Title:          eth,
			NativeAsset:    "ETH",
			BlockchainName: "Ethereum",
		},
		{
			IsActive:       &isActive,
			Title:          bsc,
			NativeAsset:    "BSC",
			BlockchainName: "Binance",
		},
	}

	var db = database.GetInstance()
	for _, blockchain := range blockchains {
		db.GetClient().FirstOrCreate(&blockchain, model.Blockchain{BlockchainName: blockchain.BlockchainName, NativeAsset: blockchain.NativeAsset})
	}
	log.Println("Blockchain Seeder executed successfully.")
}
