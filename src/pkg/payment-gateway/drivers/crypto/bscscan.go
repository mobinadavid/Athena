package crypto

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

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/blockchain-explorer")
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
		SetTimeout(time.Duration(requestTimeout) * time.Second)

	return bscscan, nil
}

func (b *Bscscan) FetchTransactions(walletAddress string) ([]models.Response, error) {
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

	if bscScanResponse.Status != "1" {
		return nil, fmt.Errorf("API error: %s", bscScanResponse.Message)
	}

	// Parse transactions

	var transactions []models.Response

	for _, tx := range bscScanResponse.Result {
		response := models.Response{
			BlockNumber:   tx["blockNumber"].(string),
			Hash:          tx["hash"].(string),
			Timestamp:     tx["timeStamp"].(string),
			From:          tx["from"].(string),
			To:            tx["to"].(string),
			Gas:           tx["gas"].(string),
			GasPrice:      tx["gasPrice"].(string),
			Confirmations: tx["confirmations"].(string),
			Amount:        tx["value"].(string),
			BlockChain:    "BSC",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
