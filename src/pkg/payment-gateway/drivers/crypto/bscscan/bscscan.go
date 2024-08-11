package bscscan

import (
	"athena/src/config"
	"athena/src/models"
	"athena/src/pkg/vault"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type Bscscan struct {
	apiClient *resty.Client
	apiKey    string
	baseUrl   string
}

func NewBscscan(baseUrl string) (*Bscscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/api-key/blockchain-explorer")
	if err != nil {
		return nil, err
	}

	bscscan := &Bscscan{
		apiClient: resty.New(),
		apiKey:    secrets.Data["bscApiKey"].(string),
		baseUrl:   baseUrl,
	}

	bscscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).SetRetryCount(maxRetry).SetRetryWaitTime(1 * time.Second)

	return bscscan, nil
}

func (b *Bscscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", b.baseUrl, walletAddress, b.apiKey)

	resp, err := b.apiClient.R().
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Bscscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var bscScanResponse models.BscScanResponse
	err = json.Unmarshal(resp.Body(), &bscScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var transactions []*models.Transaction

	for _, tx := range bscScanResponse.Result {

		var toAddresses []string
		to := fmt.Sprintf("%v", tx["to"])
		toAddresses = append(toAddresses, to)

		var timestampInt int64
		if timestamp, ok := tx["timeStamp"].(string); ok {
			// Convert timestamp to int64
			var err error
			timestampInt, err = strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing timestamp: %w", err)
			}
		} else {
			return nil, fmt.Errorf("invalid timestamp format")
		}

		// Convert Unix timestamp to UTC time
		utcTime := time.Unix(timestampInt, 0).UTC()

		// Format into desired format
		timestampStr := utcTime.Format("Jan-02-2006 03:04:05 PM UTC")

		response := &models.Transaction{
			BlockNumber:   tx["blockNumber"].(string),
			Hash:          tx["hash"].(string),
			Timestamp:     timestampStr,
			From:          tx["from"].(string),
			ToAddresses:   toAddresses,
			GasUsed:       tx["gasUsed"].(string),
			GasPrice:      tx["gasPrice"].(string),
			Confirmations: tx["confirmations"].(string),
			Amount:        tx["value"].(string),
			BlockChain:    "BSC",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
