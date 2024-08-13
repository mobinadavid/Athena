package tronscan

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
	maxRetry, _ := strconv.Atoi(configs.Get("GET_TRANSACTIONS_MAX_RETRY"))

	tronscan := &Tronscan{
		apiClient: resty.New(),
		baseUrl:   baseUrl,
	}

	tronscan.apiClient.
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(requestTimeout) * time.Second).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(5 * time.Second)

	return tronscan, nil
}

func (t *Tronscan) FetchTransactions(walletAddress string) ([]*models.Transaction, error) {
	url := fmt.Sprintf("%s/api/transaction?start=%d&limit=%d&address=%s", t.baseUrl, 0, 100, walletAddress)
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

	return parseTransactions(tronScanResponse.Data), nil
}

// Helper function to parse transactions from the Tronscan response.
func parseTransactions(txResults []map[string]interface{}) []*models.Transaction {
	var transactions []*models.Transaction

	for _, tx := range txResults {
		toAddresses := []string{fmt.Sprintf("%v", tx["toAddress"])}

		timestampStr, err := parseTimestamp(tx["timestamp"])
		if err != nil {
			// Log the error but continue processing other transactions
			fmt.Printf("warning: %v\n", err)
		}

		feeStr := parseFee(tx["cost"])
		blockStr := parseBlock(tx["block"])

		isConfirmed, ok := tx["confirmed"].(bool)
		if !ok {
			isConfirmed = false // default value if type assertion fails
		}

		transaction := &models.Transaction{
			BlockNumber: blockStr,
			Hash:        fmt.Sprintf("%v", tx["hash"]),
			Timestamp:   timestampStr,
			From:        fmt.Sprintf("%v", tx["ownerAddress"]),
			ToAddresses: toAddresses,
			Amount:      fmt.Sprintf("%v", tx["amount"]),
			Fee:         feeStr,
			IsConfirmed: isConfirmed,
			BlockChain:  "TRX",
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}

// Helper function to parse and format the timestamp.
func parseTimestamp(timestamp interface{}) (string, error) {
	timestampMs, ok := timestamp.(float64)
	if !ok {
		return "", fmt.Errorf("invalid timestamp format")
	}

	timestampSecs := int64(timestampMs / 1000)
	utcTime := time.Unix(timestampSecs, 0).UTC()
	return utcTime.Format("Jan-02-2006 03:04:05 PM UTC"), nil
}

// Helper function to parse and format the fee.
func parseFee(cost interface{}) string {
	costMap, ok := cost.(map[string]interface{})
	if !ok {
		return fmt.Sprintf("%v", cost)
	}

	fee, ok := costMap["fee"].(float64)
	if !ok {
		return fmt.Sprintf("%v", costMap["fee"])
	}
	return fmt.Sprintf("%.0f", fee)
}

// Helper function to parse and format the block number.
func parseBlock(block interface{}) string {
	switch v := block.(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%v", block)
	}
}
