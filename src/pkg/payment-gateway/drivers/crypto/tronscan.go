package crypto

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
	page      uint
	limit     uint
	baseUrl   string
}

func NewTronscan(baseUrl string, page, limit uint) (*Tronscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))

	tronscan := &Tronscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
		page:      page,
		limit:     limit,
	}

	tronscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second)

	return tronscan, nil
}

func (t *Tronscan) FetchTransactions(walletAddress string) (int64, []models.Response, error) {
	start := (t.page - 1) * t.limit
	url := fmt.Sprintf("%s/api/transaction?start=%d&limit=%d&address=%s", t.baseUrl, start, t.limit, walletAddress)
	resp, err := t.apiClient.R().
		Get(url)

	if err != nil {
		return 0, nil, fmt.Errorf("error making request to Tronscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return 0, nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var tronScanResponse models.TronScanResponse
	err = json.Unmarshal(resp.Body(), &tronScanResponse)
	if err != nil {
		return 0, nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	totalItems := int64(tronScanResponse.Total)
	// Convert TronScanResponse to your Response type
	var transactions []models.Response
	for _, tx := range tronScanResponse.Data {
		response := models.Response{
			BlockNumber: fmt.Sprintf("%v", tx["block"]),
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   fmt.Sprintf("%v", tx["timestamp"]),
			From:        fmt.Sprintf("%v", tx["ownerAddress"]),
			To:          fmt.Sprintf("%v", tx["toAddress"]),
			Gas:         fmt.Sprintf("%v", tx["cost"].(map[string]interface{})["energy_usage"]),
			Amount:      fmt.Sprintf("%v", tx["amount"]),
			Fee:         fmt.Sprintf("%v", tx["fee"]),
			IsConfirmed: tx["confirmed"].(bool),
			BlockChain:  "TRX",
		}
		transactions = append(transactions, response)
	}

	return totalItems, transactions, nil
}
