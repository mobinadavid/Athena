package fiat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const (
	baseUrl              = "https://sep.shaparak.ir"
	requestPaymentPath   = "/onlinepg/onlinepg"
	verifyPaymentPath    = "/verifyTxnRandomSessionkey/ipg/VerifyTransaction"
	requestRetryCount    = 3
	requestRetryWaitTime = 5 * time.Second
	requestTimeout       = 5 * time.Second
)

type Sep struct {
	IPG
}

type requestPaymentResponse struct {
	Status           int    `json:"status"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDesc"`
	Token            string `json:"token"`
}

type transactionDetail struct {
	RRN             string `json:"RRN"`
	RefNum          string `json:"RefNum"`
	MaskedPan       string `json:"MaskedPan"`
	HashedPan       string `json:"HashedPan"`
	TerminalNumber  int    `json:"TerminalNumber"`
	OriginalAmount  int    `json:"OrginalAmount"`
	AffectiveAmount int    `json:"AffectiveAmount"`
	StraceDate      string `json:"StraceDate"`
	StraceNo        string `json:"StraceNo"`
}

type verifyPaymentResponse struct {
	TransactionDetail transactionDetail `json:"TransactionDetail"`
	ResultCode        int               `json:"ResultCode"`
	ResultDescription string            `json:"ResultDescription"`
	Success           bool              `json:"Success"`
}

func NewSep(apiKey string) (*Sep, error) {
	client := resty.New()
	client.SetBaseURL(baseUrl).SetHeaders(map[string]string{
		"content-type": "application/json",
		"accept":       "application/json",
	}).SetRetryCount(requestRetryCount).
		SetRetryWaitTime(requestRetryWaitTime).
		SetTimeout(requestTimeout)

	sep := &Sep{}
	sep.driver = sep
	sep.SetApiKey(apiKey).SetApiClient(client).SetBaseUrl(baseUrl)

	return sep, nil
}

func (sep *Sep) GetDriverName() string {
	return "sep"
}

func (sep *Sep) RequestPayment() (string, error) {
	payload := map[string]interface{}{
		"action":      "token",
		"TerminalId":  sep.apiKey,
		"Amount":      sep.amount,
		"ResNum":      sep.internalReferenceNumber,
		"RedirectUrl": sep.callbackUrl,
		"CellNumber":  sep.mobile,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := sep.apiClient.R().SetBody(payloadJson).Post(requestPaymentPath)
	if err != nil {
		return "", err
	}

	if req.StatusCode() != http.StatusOK {
		return "", errors.New(req.Status())
	}

	var result requestPaymentResponse
	err = json.Unmarshal(req.Body(), &result)
	if err != nil {
		return "", err
	}

	if result.Status != 1 {
		return "", errors.New(result.ErrorDescription)
	}

	return fmt.Sprintf(
		"%s/%s?token=%s",
		sep.baseUrl,
		"onlinepg/sendtoken",
		result.Token,
	), nil
}

func (sep *Sep) VerifyPayment() (interface{}, error) {
	payload := map[string]interface{}{
		"RefNum":         sep.paymentVerificationToken,
		"TerminalNumber": sep.apiKey,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := sep.apiClient.R().SetBody(payloadJson).Post(verifyPaymentPath)
	if err != nil {
		return "", err
	}

	if req.StatusCode() != http.StatusOK {
		return "", errors.New(req.Status())
	}

	var result verifyPaymentResponse
	err = json.Unmarshal(req.Body(), &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
