package models

type Response struct {
	BlockNumber   string   `json:"blockNumber"`
	Hash          string   `json:"hash"`
	Timestamp     string   `json:"timeStamp"`
	From          string   `json:"from"`
	ToAddresses   []string `json:"to"`
	Gas           string   `json:"gas"`
	GasPrice      string   `json:"gasPrice"`
	Fee           string   `json:"fee"`
	IsConfirmed   bool     `json:"is_confirmed"`
	Confirmations string   `json:"confirmations"`
	Amount        string   `json:"amount"`
	BlockChain    string   `json:"blockChain"`
}
