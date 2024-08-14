package models

type TronScanResponse struct {
	Total int                      `json:"total"`
	Data  []map[string]interface{} `json:"data"`
}
