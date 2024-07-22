package tronscan

import (
	"athena/src/services/payment-gateway/drivers/crypto/tronscan/model"
	transaction_response "athena/src/services/transaction-response"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func FetchTronTransactions(walletAddress string, baseUrl string) ([]transaction_response.Response, error) {
	client := resty.New()

	url := fmt.Sprintf("%s/api/transaction?limit=20&start=0&address=%s", baseUrl, walletAddress)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Tronscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var tronScanResponse model.TronScanResponse
	err = json.Unmarshal(resp.Body(), &tronScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Convert TronScanResponse to your Response type
	var transactions []transaction_response.Response
	for _, tx := range tronScanResponse.Records {
		response := transaction_response.Response{
			BlockNumber: fmt.Sprintf("%v", tx["block"]),
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   fmt.Sprintf("%v", tx["timestamp"]),
			From:        fmt.Sprintf("%v", tx["ownerAddress"]),
			To:          fmt.Sprintf("%v", tx["toAddress"]),
			Gas:         fmt.Sprintf("%v", tx["cost"].(map[string]interface{})["energy_usage"]),
			Amount:      fmt.Sprintf("%v", tx["amount"]),
			IsConfirmed: tx["confirmed"].(bool),
			BlockChain:  "TRX",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
