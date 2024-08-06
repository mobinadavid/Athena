package btcscan

import (
	"athena/src/config"
	"athena/src/models"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type Btcscan struct {
	apiClient *resty.Client
	baseUrl   string
}

func NewBtcscan(baseUrl string) (*Btcscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	btcscan := &Btcscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
	}

	btcscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).SetRetryCount(maxRetry).SetRetryWaitTime(1 * time.Second)

	return btcscan, nil
}

func (t *Btcscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/rawaddr/%s", t.baseUrl, walletAddress)
	resp, err := t.apiClient.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error making request to Btcscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var btcScanResponse models.BtcScanResponse
	err = json.Unmarshal(resp.Body(), &btcScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Convert TronScanResponse to your Transaction type
	var transactions []*models.Transaction
	for _, tx := range btcScanResponse.Txs {
		inputs := tx["inputs"].([]interface{})
		var fromAddr string

		if len(inputs) > 0 {
			input := inputs[0].(map[string]interface{})
			prevOut := input["prev_out"].(map[string]interface{})
			fromAddr = fmt.Sprintf("%v", prevOut["addr"])
		}

		outAddresses := tx["out"].([]interface{})
		var toAddresses []string
		for _, out := range outAddresses {
			outMap := out.(map[string]interface{})
			if addr, ok := outMap["addr"].(string); ok {
				toAddresses = append(toAddresses, addr)
			}
		}

		var timestampInt int64
		if timestamp, ok := tx["time"].(float64); ok {
			timestampInt = int64(timestamp)
		} else {
			return nil, fmt.Errorf("invalid timestamp format")
		}

		// Format timestamp to a readable date/time string
		timestampStr := time.Unix(timestampInt, 0).Format(time.RFC3339)

		response := &models.Transaction{
			BlockNumber: fmt.Sprintf("%v", tx["block_height"]),
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   timestampStr,
			From:        fromAddr,
			ToAddresses: toAddresses,
			Amount:      fmt.Sprintf("%v", tx["result"]),
			Fee:         fmt.Sprintf("%v", tx["fee"]),
			BlockChain:  "BTC",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
