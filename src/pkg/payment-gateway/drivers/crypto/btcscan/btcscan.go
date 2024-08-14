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
	baseUrl   string
}

func NewBtcscan(baseUrl string) (*Btcscan, error) {
	configs := config.GetInstance()
	requestTimeout, _ := strconv.Atoi(configs.Get("EXPLORER_REQUEST_TIMEOUT"))
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	btcscan := &Btcscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
	}

	btcscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(5 * time.Second)

	return btcscan, nil
}

func (t *Btcscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/rawaddr/%s", t.baseUrl, walletAddress)
	resp, err := t.apiClient.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error making request to Btcscan: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	var btcScanResponse models.BtcScanResponse
	err = json.Unmarshal(resp.Body(), &btcScanResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return parseTransactions(btcScanResponse.Txs)
}

// Helper function to parse transactions from the BtcScan response.
func parseTransactions(txs []map[string]interface{}) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	for _, tx := range txs {
		fromAddr, err := parseFromAddress(tx["inputs"])
		if err != nil {
			return nil, err
		}

		toAddresses, err := parseToAddresses(tx["out"])
		if err != nil {
			return nil, err
		}

		timestampStr, err := parseTimestamp(tx["time"])
		if err != nil {
			return nil, err
		}

		transaction := &models.Transaction{
			BlockNumber: fmt.Sprintf("%v", tx["block_height"]),
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   timestampStr,
			From:        fromAddr,
			ToAddresses: toAddresses,
			Amount:      fmt.Sprintf("%v", tx["result"]),
			Fee:         fmt.Sprintf("%v", tx["fee"]),
			BlockChain:  "BTC",
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// Helper function to parse the "from" address from transaction inputs.
func parseFromAddress(inputs interface{}) (string, error) {
	if inputs == nil {
		return "", nil
	}

	inputsSlice, ok := inputs.([]interface{})
	if !ok || len(inputsSlice) == 0 {
		return "", fmt.Errorf("invalid or empty 'inputs'")
	}

	input, ok := inputsSlice[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid 'input' format")
	}

	prevOut, ok := input["prev_out"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid 'prev_out' format")
	}

	addr, ok := prevOut["addr"].(string)
	if !ok {
		return "", fmt.Errorf("invalid 'addr' format")
	}
	return addr, nil
}

// Helper function to parse "to" addresses from transaction outputs.
func parseToAddresses(out interface{}) ([]string, error) {
	if out == nil {
		return nil, nil
	}

	outSlice, ok := out.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid 'out' format")
	}

	var toAddresses []string
	for _, item := range outSlice {
		outMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid 'out' item format")
		}
		addr, ok := outMap["addr"].(string)
		if ok {
			toAddresses = append(toAddresses, addr)
		}
	}
	return toAddresses, nil
}

// Helper function to parse and format the timestamp.
func parseTimestamp(timestamp interface{}) (string, error) {
	timestampFloat, ok := timestamp.(float64)
	if !ok {
		return "", fmt.Errorf("invalid timestamp format")
	}

	timestampInt := int64(timestampFloat)
	utcTime := time.Unix(timestampInt, 0).UTC()
	return utcTime.Format("Jan-02-2006 03:04:05 PM UTC"), nil
}
