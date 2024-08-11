package services

import (
	"athena/src/api/http/requests"
	payment_gateway "athena/src/pkg/payment-gateway"
	"github.com/google/uuid"
)

type IIpgService interface {
	RequestPayment(request *requests.PaymentRequest) (string, error)
}
type IpgService struct {
}

func (s *IpgService) RequestPayment(request *requests.PaymentRequest) (string, error) {
	ipgDriver, err := payment_gateway.NewIPG()
	if err != nil {
		return "", err
	}

	redirectUrl, err := ipgDriver.SetAmount(request.Amount).
		SetInternalReferenceNumber(uuid.NewString()).
		SetCallbackUrl("https://").
		SetMobile(request.Mobile).
		RequestPayment()

	if err != nil {
		return "", err
	}
	return redirectUrl, nil
}
