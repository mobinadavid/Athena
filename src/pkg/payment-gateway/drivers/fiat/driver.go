package fiat

import (
	"github.com/go-resty/resty/v2"
)

type IIPGDriver interface {
	SetBaseUrl(string) IIPGDriver
	SetApiKey(string) IIPGDriver
	GetDriverName() string
	SetMobile(string) IIPGDriver
	SetDescription(string) IIPGDriver
	SetAmount(float64) IIPGDriver
	SetCallbackUrl(string) IIPGDriver
	SetInternalReferenceNumber(string) IIPGDriver
	SetPaymentVerificationToken(string) IIPGDriver
	SetApiClient(*resty.Client) IIPGDriver
	RequestPayment() (string, error)
	VerifyPayment() (interface{}, error)
}

type IPG struct {
	driver                   IIPGDriver
	apiClient                *resty.Client
	baseUrl                  string
	apiKey                   string
	mobile                   string
	description              string
	callbackUrl              string
	internalReferenceNumber  string
	paymentVerificationToken string
	amount                   float64
}

func (ipg *IPG) GetDriverName() string {
	return ipg.driver.GetDriverName()
}

func (ipg *IPG) RequestPayment() (string, error) {
	return ipg.driver.RequestPayment()
}

func (ipg *IPG) VerifyPayment() (interface{}, error) {
	return ipg.driver.VerifyPayment()
}

func (ipg *IPG) SetBaseUrl(url string) IIPGDriver {
	ipg.baseUrl = url
	return ipg
}

func (ipg *IPG) SetApiKey(key string) IIPGDriver {
	ipg.apiKey = key
	return ipg
}

func (ipg *IPG) SetCallbackUrl(callbackUrl string) IIPGDriver {
	ipg.callbackUrl = callbackUrl
	return ipg
}

func (ipg *IPG) SetInternalReferenceNumber(internalReferenceNumber string) IIPGDriver {
	ipg.internalReferenceNumber = internalReferenceNumber
	return ipg
}

func (ipg *IPG) SetMobile(mobile string) IIPGDriver {
	ipg.mobile = mobile
	return ipg
}

func (ipg *IPG) SetDescription(description string) IIPGDriver {
	ipg.description = description
	return ipg
}

func (ipg *IPG) SetAmount(amount float64) IIPGDriver {
	ipg.amount = amount
	return ipg
}

func (ipg *IPG) SetPaymentVerificationToken(token string) IIPGDriver {
	ipg.paymentVerificationToken = token
	return ipg
}

func (ipg *IPG) SetApiClient(client *resty.Client) IIPGDriver {
	ipg.apiClient = client
	return ipg
}
