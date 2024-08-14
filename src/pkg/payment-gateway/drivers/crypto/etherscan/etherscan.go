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
		SetTimeout(time.Duration(requestTimeout) * time.Second).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(5 * time.Second)

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

	return parseTransactions(etherScanResponse.Result), nil
}

// Helper function to parse transactions from the Etherscan response.
func parseTransactions(txResults []map[string]interface{}) []*models.Transaction {
	var transactions []*models.Transaction

	for _, tx := range txResults {
		toAddresses := []string{fmt.Sprintf("%v", tx["to"])}

		timestampStr, err := parseTimestamp(tx["timeStamp"])
		if err != nil {
			// Log the error but continue processing other transactions
			fmt.Printf("warning: %v\n", err)
		}

		transaction := &models.Transaction{
			BlockNumber:   fmt.Sprintf("%v", tx["blockNumber"]),
			Hash:          fmt.Sprintf("%v", tx["hash"]),
			Timestamp:     timestampStr,
			From:          fmt.Sprintf("%v", tx["from"]),
			ToAddresses:   toAddresses,
			GasUsed:       fmt.Sprintf("%v", tx["gasUsed"]),
			GasPrice:      fmt.Sprintf("%v", tx["gasPrice"]),
			Confirmations: fmt.Sprintf("%v", tx["confirmations"]),
			Amount:        fmt.Sprintf("%v", tx["value"]),
			BlockChain:    "ETH",
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}

// Helper function to parse and format the timestamp.
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
