package tronscan

import (
	"athena/src/config"
	"athena/src/models"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type Tronscan struct {
	apiClient *resty.Client
	baseUrl   string
}

func NewTronscan(baseUrl string) (*Tronscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	tronscan := &Tronscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
	}

	tronscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).SetRetryCount(maxRetry).SetRetryWaitTime(1 * time.Second)

	return tronscan, nil
}

func (t *Tronscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/api/transaction?start=%d&limit=%d&address=%s", t.baseUrl, 0, 100, walletAddress)
	resp, err := t.apiClient.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error making request to Tronscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var tronScanResponse models.TronScanResponse
	err = json.Unmarshal(resp.Body(), &tronScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Convert TronScanResponse to your Transaction type
	var transactions []*models.Transaction
	for _, tx := range tronScanResponse.Data {

		var toAddresses []string
		to := fmt.Sprintf("%v", tx["toAddress"])
		toAddresses = append(toAddresses, to)

		var timestampStr string
		if timestampMs, ok := tx["timestamp"].(float64); ok {
			// Convert milliseconds to seconds
			timestampSecs := int64(timestampMs / 1000)
			// Convert Unix timestamp to time.Time
			utcTime := time.Unix(timestampSecs, 0).UTC()
			// Format time.Time to a string
			timestampStr = utcTime.Format("Jan-02-2006 03:04:05 PM UTC")
		} else {
			// Handle the case where timestamp is not a float64
			timestampStr = fmt.Sprintf("%v", tx["timestamp"])
		}

		var feeStr string
		if cost, ok := tx["cost"].(map[string]interface{}); ok {
			if fee, ok := cost["fee"].(float64); ok {
				feeStr = fmt.Sprintf("%.0f", fee) // Adjust precision as needed
			} else {
				feeStr = fmt.Sprintf("%v", cost["fee"])
			}
		}

		var blockstr string
		if block, ok := tx["block"].(float64); ok {
			blockstr = fmt.Sprintf("%.0f", block)
		} else {
			blockstr = fmt.Sprintf("%v", tx["block"])
		}

		response := &models.Transaction{
			BlockNumber: blockstr,
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   timestampStr,
			From:        fmt.Sprintf("%v", tx["ownerAddress"]),
			ToAddresses: toAddresses,
			Amount:      fmt.Sprintf("%v", tx["amount"]),
			Fee:         feeStr,
			IsConfirmed: tx["confirmed"].(bool),
			BlockChain:  "TRX",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
