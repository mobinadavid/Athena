package models

import "time"

type TransactionResponse struct {
	Hash          string    `json:"hash"`
	IsConfirmed   bool      `json:"is_confirmed"`
	BlockNumber   int64     `json:"block_number"`
	Timestamp     time.Time `json:"timeStamp"`
	From          string    `json:"from"`
	ToAddresses   []string  `json:"to"`
	Value         float64   `json:"value"`
	Fee           float64   `json:"transaction_fee"`
	BlockChain    string    `json:"blockChain"`
	Confirmations int       `json:"confirmations"`
}
