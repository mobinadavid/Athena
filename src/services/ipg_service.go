package services

import (
	"athena/src/api/http/requests"
	"athena/src/config"
	"athena/src/models"
	payment_gateway "athena/src/pkg/payment-gateway"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type IIpgService interface {
	RequestPayment(request *requests.PaymentRequest) (interface{}, error)
	VerifyPayment(uuid *uuid.UUID, tx requests.IpgCallbackRequest) (interface{}, error)
	GetIGPByUuid(uuid *uuid.UUID) (*models.IGPModel, error)
}
type IpgService struct {
	IIGPService IIGPService
}

func (service *IpgService) RequestPayment(request *requests.PaymentRequest) (interface{}, error) {
	// IPG
	ipgDriver, err := payment_gateway.NewIPG(request.Driver)
	if err != nil {
		return "", err
	}

	// IGP
	igp, err := service.IIGPService.Create(&models.IGPModel{
		Ipg:         ipgDriver.GetDriverName(),
		Amount:      request.Amount,
		CallbackUrl: request.CallBackUrl,
		Status:      "pending",
	})

	if err != nil {
		return nil, err
	}

	callbackUrl := fmt.Sprintf("%s://%s/%s/%s",
		"https",
		config.GetInstance().Get("APP_HOST"),
		"api/v1/ipg/ipg-callback",
		igp.Uuid.String(),
	)

	ipgPaymentUrl, err := ipgDriver.SetAmount(request.Amount).
		SetMobile(request.Mobile).
		SetInternalReferenceNumber(igp.Uuid.String()).
		SetCallbackUrl(callbackUrl).
		RequestPayment()
	if err != nil {
		return nil, err
	}

	return ipgPaymentUrl, nil
}

func (service *IpgService) VerifyPayment(uuid *uuid.UUID, tx requests.IpgCallbackRequest) (interface{}, error) {
	//Get IGP
	igp, err := service.IIGPService.GetForVerification(uuid, tx.RefNum)
	if err != nil {
		return nil, err
	}

	// initialize IPG
	ipgDriver, err := payment_gateway.NewIPG(igp.Ipg)
	if err != nil {
		return nil, err
	}

	// Verify Payment
	verification, err := ipgDriver.
		SetPaymentVerificationToken(tx.RefNum).
		VerifyPayment()
	if err != nil {
		igp.Status = "failed"
		txJson, _ := json.Marshal(tx)
		igp.Receipt = txJson
		_, err = service.IIGPService.Update(igp)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	igp.Status = "done"
	igp.IssuerReferenceNumber = tx.RefNum
	txJson, _ := json.Marshal(tx)
	igp.Receipt = txJson
	// Update IGP
	_, err = service.IIGPService.Update(igp)
	if err != nil {
		return nil, err
	}

	return verification, nil
}

func (service *IpgService) GetIGPByUuid(uuid *uuid.UUID) (*models.IGPModel, error) {
	return service.IIGPService.GetByUuid(uuid)
}
