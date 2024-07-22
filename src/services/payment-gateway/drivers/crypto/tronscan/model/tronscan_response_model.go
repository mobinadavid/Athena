package model

type TronScanResponse struct {
	Total   int                      `json:"total"`
	Records []map[string]interface{} `json:"data"`
}
