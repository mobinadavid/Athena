package etherscan

import (
	"athena/src/config"
	"athena/src/pkg/vault"
	"athena/src/services/payment-gateway/drivers/crypto/etherscan/model"
	"athena/src/services/wallet-address/transaction-response"
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
	requestTimeout, _ := strconv.Atoi(configs.Get("ETH_REQUEST_TIMEOUT"))

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), configs.Get("APP_NAME")+"/blockchain-explorer")
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
		SetTimeout(time.Duration(requestTimeout) * time.Second)

	if proxy := configs.Get("ETH_PROXY"); proxy != "" {
		etherscan.apiClient.SetProxy(proxy)
	}

	return etherscan, nil
}

func (e *Etherscan) FetchEthTransaction(walletAddress string) ([]transaction_response.Response, error) {
	url := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&apikey=%s", e.baseUrl, walletAddress, e.apiKey)

	resp, err := e.apiClient.R().
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Etherscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var etherScanResponse model.EtherScanResponse
	err = json.Unmarshal(resp.Body(), &etherScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if etherScanResponse.Status != "1" {
		return nil, fmt.Errorf("API error: %s", etherScanResponse.Message)
	}
	// Parse transactions

	var transactions []transaction_response.Response
	for _, tx := range etherScanResponse.Result {

		response := transaction_response.Response{
			BlockNumber:   tx["blockNumber"].(string),
			Hash:          tx["hash"].(string),
			Timestamp:     tx["timeStamp"].(string),
			From:          tx["from"].(string),
			To:            tx["to"].(string),
			Gas:           tx["gas"].(string),
			GasPrice:      tx["gasPrice"].(string),
			Confirmations: tx["confirmations"].(string),
			Amount:        tx["value"].(string),
			BlockChain:    "ETH",
		}
		transactions = append(transactions, response)
	}

	return transactions, nil
}
