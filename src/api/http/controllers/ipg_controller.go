package controllers

import (
	"athena/src/api/http/requests"
	"athena/src/api/http/response"
	"athena/src/pkg/i18n"
	"athena/src/pkg/validator"
	"athena/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	redirectUrl, err := controller.IIpgService.RequestPayment(&req)
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
			"redirect_url": redirectUrl,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}

func (controller *IpgController) VerifyPayment(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")
	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)

	igp, err := controller.IIpgService.GetIGPByUuid(&id)
	if err != nil {
		response.Api(c).
			SetMessage(err.Error()).
			Send()
		return
	}

	clientUrl := igp.CallbackUrl + "/"

	// todo: tx is for SEP right now, should be dynamic for multiple ipgs.
	var tx requests.IpgCallbackRequest
	if err = c.ShouldBind(&tx); err != nil {
		c.Redirect(http.StatusMovedPermanently, clientUrl+uuidStr)
		return
	}

	_, err = controller.IIpgService.VerifyPayment(&id, tx)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, clientUrl+uuidStr)
		return
	}

	c.Redirect(http.StatusMovedPermanently, clientUrl+uuidStr)
	return
}

func (controller *IpgController) GetIGPByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	igp, err := controller.IIpgService.GetIGPByUuid(&id)
	if err != nil {
		response.Api(c).
			SetStatusCode(http.StatusNotFound).
			SetMessage(err.Error()).
			Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"igp": igp,
		}).
		Send()
	return
}
