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
	//get api key from our vault
	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/blockchain-explorer/api-key")
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
		SetTimeout(time.Duration(requestTimeout) * time.Second).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(5 * time.Second)

	return bscscan, nil
}

// FetchTransactions retrieves transactions for a given wallet address from Bscscan.
func (b *Bscscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", b.baseUrl, walletAddress, b.apiKey)

	resp, err := b.apiClient.R().Get(url)
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

	return parseTransactions(bscScanResponse.Result)
}

// Helper function to parse transactions from the BscScan response.
func parseTransactions(results []map[string]interface{}) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	for _, tx := range results {
		var toAddresses []string
		to := fmt.Sprintf("%v", tx["to"])
		toAddresses = append(toAddresses, to)

		timestampStr, err := parseTimestamp(tx["timeStamp"])
		if err != nil {
			return nil, err
		}

		transaction := &models.Transaction{
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
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// Helper function to parse timestamp from the transaction.
func parseTimestamp(timestamp interface{}) (string, error) {
	timestampStr, ok := timestamp.(string)
	if !ok {
		return "", fmt.Errorf("invalid timestamp format")
	}

	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing timestamp: %w", err)
	}

	utcTime := time.Unix(timestampInt, 0).UTC()
	return utcTime.Format("Jan-02-2006 03:04:05 PM UTC"), nil
}
