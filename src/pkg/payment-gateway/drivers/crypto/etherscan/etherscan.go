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
	page      uint
	limit     uint
}

func NewEtherscan(baseUrl string, page, limit uint) (*Etherscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/blockchain-explorer")
	if err != nil {
		return nil, err
	}

	etherscan := &Etherscan{
		apiClient: resty.New(),
		apiKey:    secrets.Data["ethApiKey"].(string),
		baseUrl:   baseUrl,
		page:      page,
		limit:     limit,
	}

	etherscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).SetRetryCount(maxRetry).SetRetryWaitTime(1 * time.Second)

	return etherscan, nil
}

func (e *Etherscan) FetchTransactions(walletAddress string) (int64, []*models.Transaction, error) {
	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", e.baseUrl, walletAddress, e.apiKey)

	resp, err := e.apiClient.R().
		Get(url)
	if err != nil {
		return 0, nil, fmt.Errorf("error making request to Etherscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return 0, nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var etherScanResponse models.EtherScanResponse
	err = json.Unmarshal(resp.Body(), &etherScanResponse)
	if err != nil {
		return 0, nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if etherScanResponse.Status != "1" {
		return 0, nil, fmt.Errorf("API error: %s", etherScanResponse.Message)
	}
	// Parse transactions
	startIndex := (e.page - 1) * e.limit
	endIndex := e.page * e.limit
	totalItems := int64(len(etherScanResponse.Result))

	// Ensure startIndex and endIndex are within bounds
	if startIndex > uint(len(etherScanResponse.Result)) {
		startIndex = uint(len(etherScanResponse.Result))
	}
	if endIndex > uint(len(etherScanResponse.Result)) {
		endIndex = uint(len(etherScanResponse.Result))
	}
	if startIndex == 0 && endIndex == 0 {
		endIndex = uint(len(etherScanResponse.Result))
	}

	// Get the subset of data for the requested page
	paginatedData := etherScanResponse.Result[startIndex:endIndex]

	var transactions []*models.Transaction
	for _, tx := range paginatedData {

		var toAddresses []string
		to := fmt.Sprintf("%v", tx["to"])
		toAddresses = append(toAddresses, to)

		response := &models.Transaction{
			BlockNumber:   tx["blockNumber"].(string),
			Hash:          tx["hash"].(string),
			Timestamp:     tx["timeStamp"].(string),
			From:          tx["from"].(string),
			ToAddresses:   toAddresses,
			Gas:           tx["gas"].(string),
			GasPrice:      tx["gasPrice"].(string),
			Confirmations: tx["confirmations"].(string),
			Amount:        tx["value"].(string),
			BlockChain:    "ETH",
		}
		transactions = append(transactions, response)
	}

	return totalItems, transactions, nil
}
