package fiat

import "golang.org/x/text/currency"

type IFiatDriver interface {
	GenerateTransaction()
	VerifyTransaction()
}

type FiatDriver struct {
	BaseUrl     string
	ApiKey      string
	CallbackUrl string
	Mobile      string
	Amount      currency.Amount
	Description string
}
