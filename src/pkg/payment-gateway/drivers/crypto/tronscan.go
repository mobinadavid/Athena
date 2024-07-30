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
	baseUrl   string
}

func NewTronscan(baseUrl string) (*Tronscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))

	tronscan := &Tronscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
	}

	tronscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second)

	return tronscan, nil
}

func (t *Tronscan) FetchTransactions(walletAddress string) ([]models.Response, error) {
	url := fmt.Sprintf("%s/api/transaction?limit=20&start=0&address=%s", t.baseUrl, walletAddress)

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

	// Convert TronScanResponse to your Response type
	var transactions []models.Response
	for _, tx := range tronScanResponse.Records {
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

	return transactions, nil
}
