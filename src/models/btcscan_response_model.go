package models

type BtcScanResponse struct {
	TotalTransactions int                      `json:"n_tx"`
	Txs               []map[string]interface{} `json:"txs"`
}
