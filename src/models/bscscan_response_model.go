package models

type BscScanResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Result  []map[string]interface{} `json:"result"`
}
