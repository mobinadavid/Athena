package controllers

import (
	"athena/src/api/http/requests"
	"athena/src/api/http/response"
	"athena/src/database/scopes"
	"athena/src/pkg/i18n"
	"athena/src/pkg/validator"
	"athena/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type WalletAddressController struct {
	IWalletAddressService services.IWalletAddressService
}

func (controller *WalletAddressController) GetWalletAddress(c *gin.Context) {
	var req requests.GetWalletAddress
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

	walletAddress, err := controller.IWalletAddressService.GetWalletAddress(&req)
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
			"wallet_address": walletAddress,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}

func (controller *WalletAddressController) GetList(c *gin.Context) {
	var walletAddress *scopes.PaginateModel

	walletAddress, err := controller.IWalletAddressService.GetList()

	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"wallet_addresses": walletAddress,
		}).Send()
	return
}

func (controller *WalletAddressController) GetByUuid(c *gin.Context) {
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	walletAddress, err := controller.IWalletAddressService.GetByUuid(&id)

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
			"wallet_address:": walletAddress,
		}).
		Send()
	return
}

func (controller *WalletAddressController) Create(c *gin.Context) {
	var req requests.CreateWalletAddressRequest

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

	walletAddress, err := controller.IWalletAddressService.Create(&req)
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
			"wallet_address": walletAddress,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()
}

func (controller *WalletAddressController) Delete(c *gin.Context) {
	//Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}
	err = controller.IWalletAddressService.Delete(&id)

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
		SetData(map[string]interface{}{}).
		Send()
	return
}

func (controller *WalletAddressController) Update(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	var req requests.CreateWalletAddressRequest

	// Bind check payload.
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

	walletAddress, err := controller.IWalletAddressService.Update(&id, &req)

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
			"wallet_address": walletAddress,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}

func (controller *WalletAddressController) GetTransactions(c *gin.Context) {
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	walletAddress, err := controller.IWalletAddressService.GetByUuid(&uuid)

	if err != nil {
		response.Api(c).
			SetStatusCode(http.StatusNotFound).
			SetMessage(err.Error()).
			Send()
		return
	}

	transactions, err := controller.IWalletAddressService.GetTransactions(walletAddress)
	if err != nil {
		response.Api(c).
			SetMessage(err.Error()).
			Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"transactions": transactions,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}
