package transaction_response

type Response struct {
	BlockNumber   string `json:"blockNumber"`
	Hash          string `json:"hash"`
	Timestamp     string `json:"timeStamp"`
	BlockHash     string `json:"blockHash"`
	From          string `json:"from"`
	To            string `json:"to"`
	Gas           string `json:"gas"`
	GasPrice      string `json:"gasPrice"`
	Nonce         string `json:"nonce"`
	Confirmations string `json:"confirmations"`
	BlockChain    string `json:"blockChain"`
}
