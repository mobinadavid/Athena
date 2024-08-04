package seeders

import (
	"athena/src/database"
	"athena/src/models"
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

	titleBtc := map[string]string{
		"en": "Bitcoin",
		"fa": "بیت کوین",
	}
	btc, err := json.Marshal(titleBtc)
	if err != nil {

	}

	isActive := true
	blockchains := []*models.Blockchain{
		{
			IsActive:    &isActive,
			Title:       tron,
			NativeAsset: "TRX",
			Name:        "Tron",
		},
		{
			IsActive:    &isActive,
			Title:       eth,
			NativeAsset: "ETH",
			Name:        "Ethereum",
		},
		{
			IsActive:    &isActive,
			Title:       bsc,
			NativeAsset: "BSC",
			Name:        "Binance",
		},
		{IsActive: &isActive,
			Title:       btc,
			NativeAsset: "BTC",
			Name:        "Bitcoin",
		}}

	var db = database.GetInstance()
	for _, blockchain := range blockchains {
		db.GetClient().FirstOrCreate(&blockchain, models.Blockchain{Name: blockchain.Name, NativeAsset: blockchain.NativeAsset})
	}
	log.Println("Blockchain Seeder executed successfully.")
}
