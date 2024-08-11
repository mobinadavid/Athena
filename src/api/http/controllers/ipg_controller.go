package controllers

import (
	"athena/src/api/http/requests"
	"athena/src/api/http/response"
	"athena/src/pkg/i18n"
	"athena/src/pkg/validator"
	"athena/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IpgController struct {
	IIpgService services.IIpgService
}

func (controller *IpgController) RequestPayment(c *gin.Context) {
	var req requests.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Api(c).Send()
		return
	}

	// Validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		response.Api(c).
			SetStatusCode(http.StatusUnprocessableEntity).
			SetErrors(err).
			Send()
		return
	}

	token, err := controller.IIpgService.RequestPayment(&req)
	if err != nil {
		response.Api(c).
			SetMessage(err.Error()).
			Send()
		return
	}

	// Return response.
	response.Api(c).
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"token": token,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}
