package crypto

import (
	transaction_response "athena/src/services/transaction-response"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EtherscanResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Result  []map[string]interface{} `json:"result"`
}

func FetchEthTransaction(walletAddress string, baseUrl string, apiKey string) ([]transaction_response.Response, error) {
	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", baseUrl, walletAddress, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Etherscan: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var etherscanResponse EtherscanResponse
	err = json.Unmarshal(body, &etherscanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if etherscanResponse.Status != "1" {
		return nil, fmt.Errorf("API error: %s", etherscanResponse.Message)
	}

	var transactions []transaction_response.Response
	for _, tx := range etherscanResponse.Result {
		response := transaction_response.Response{
			BlockNumber:   tx["blockNumber"].(string),
			Hash:          tx["hash"].(string),
			Timestamp:     tx["timeStamp"].(string),
			BlockHash:     tx["blockHash"].(string),
			From:          tx["from"].(string),
			To:            tx["to"].(string),
			Gas:           tx["gas"].(string),
			GasPrice:      tx["gasPrice"].(string),
			Nonce:         tx["nonce"].(string),
			Confirmations: tx["confirmations"].(string),
			BlockChain:    "ETH",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
