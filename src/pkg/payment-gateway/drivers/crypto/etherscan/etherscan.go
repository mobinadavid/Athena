package etherscan

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

type Etherscan struct {
	apiClient *resty.Client
	apiKey    string
	baseUrl   string
}

func NewEtherscan(baseUrl string) (*Etherscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/blockchain-explorer/api-key")
	if err != nil {
		return nil, err
	}

	etherscan := &Etherscan{
		apiClient: resty.New(),
		apiKey:    secrets.Data["ethApiKey"].(string),
		baseUrl:   baseUrl,
	}

	etherscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).SetRetryCount(maxRetry).SetRetryWaitTime(1 * time.Second)

	return etherscan, nil
}

func (e *Etherscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", e.baseUrl, walletAddress, e.apiKey)

	resp, err := e.apiClient.R().
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Etherscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var etherScanResponse models.EtherScanResponse
	err = json.Unmarshal(resp.Body(), &etherScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var transactions []*models.Transaction
	for _, tx := range etherScanResponse.Result {
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
			BlockChain:    "ETH",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
