package models

type Response struct {
	Hash        string   `json:"hash"`
	IsConfirmed bool     `json:"is_confirmed"`
	BlockNumber string   `json:"blockNumber"`
	Timestamp   string   `json:"timeStamp"`
	From        string   `json:"from"`
	ToAddresses []string `json:"to"`
	Value       string   `json:"value"`
	Fee         string   `json:"transaction_fee"`
	BlockChain  string   `json:"blockChain"`
}
