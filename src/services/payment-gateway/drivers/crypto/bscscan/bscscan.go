package bscscan

import (
	"athena/src/services/payment-gateway/drivers/crypto/bscscan/model"
	transaction_response "athena/src/services/transaction-response"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func FetchBscTransaction(walletAddress string, baseUrl string, apiKey string) ([]transaction_response.Response, error) {

	// Create a new Resty client
	client := resty.New()

	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", baseUrl, walletAddress, apiKey)

	// Make the HTTP GET request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Etherscan: %w", err)
	}

	// Check for successful response status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	// Unmarshal the response body into EtherScanResponse
	var bscScanResponse model.BscScanResponse
	err = json.Unmarshal(resp.Body(), &bscScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if bscScanResponse.Status != "1" {
		return nil, fmt.Errorf("API error: %s", bscScanResponse.Message)
	}

	// Parse transactions

	var transactions []transaction_response.Response

	for _, tx := range bscScanResponse.Result {
		response := transaction_response.Response{
			BlockNumber:   tx["blockNumber"].(string),
			Hash:          tx["hash"].(string),
			Timestamp:     tx["timeStamp"].(string),
			From:          tx["from"].(string),
			To:            tx["to"].(string),
			Gas:           tx["gas"].(string),
			GasPrice:      tx["gasPrice"].(string),
			Confirmations: tx["confirmations"].(string),
			Amount:        tx["value"].(string),
			BlockChain:    "BSC",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
