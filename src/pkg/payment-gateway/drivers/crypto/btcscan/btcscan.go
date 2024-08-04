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
	page      uint
	limit     uint
	baseUrl   string
}

func NewBtcscan(baseUrl string, page, limit uint) (*Btcscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	btcscan := &Btcscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
		page:      page,
		limit:     limit,
	}

	btcscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).SetRetryCount(maxRetry).SetRetryWaitTime(1 * time.Second)

	return btcscan, nil
}

func (t *Btcscan) FetchTransactions(walletAddress string) (int64, []models.Response, error) {
	offset := (t.page - 1) * t.limit
	url := fmt.Sprintf("%s/rawaddr/%s?limit=%d&offset=%d", t.baseUrl, walletAddress, t.limit, offset)
	resp, err := t.apiClient.R().
		Get(url)

	if err != nil {
		return 0, nil, fmt.Errorf("error making request to Btcscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return 0, nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var btcScanResponse models.BtcScanResponse
	err = json.Unmarshal(resp.Body(), &btcScanResponse)
	if err != nil {
		return 0, nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	totalItems := int64(btcScanResponse.TotalTransactions)
	// Convert TronScanResponse to your Response type
	var transactions []models.Response
	for _, tx := range btcScanResponse.Txs {
		response := models.Response{
			BlockNumber: fmt.Sprintf("%v", tx["block_height"]),
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   fmt.Sprintf("%v", tx["time"]),
			//From:        fmt.Sprintf("%v", tx["inputs"].(map[string]interface{})["addr"]),
			//		To:          fmt.Sprintf("%v", tx["out"].(map[string]interface{})["addr"]),
			Amount:     fmt.Sprintf("%v", tx["result"]),
			Fee:        fmt.Sprintf("%v", tx["fee"]),
			BlockChain: "BTC",
		}
		transactions = append(transactions, response)
	}

	return totalItems, transactions, nil
}
