package controllers

import (
	"athena/src/api/http/requests"
	"athena/src/api/http/response"
	"athena/src/database/scopes"
	"athena/src/pkg/i18n"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletAddressController struct {
	IWalletAddressService services.IWalletAddressService
}

func (controller *WalletAddressController) GetList(c *gin.Context) {
	var walletAddress *scopes.PaginatedModel
	filters := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if key != "page" && key != "limit" && key != "sort_by" && key != "sort_order" && key != "created_after" && key != "created_before" {
			filters[key] = values[0]
		}
	}

	params := &scopes.QueryBuilderModel{
		Page:      uint(c.GetInt("page")),
		Limit:     uint(c.GetInt("limit")),
		SortBy:    c.GetString("sort_by"),
		SortOrder: c.GetString("sort_order"),
		Filters:   filters,
	}

	if createdAfterStr := c.Query("created_after"); createdAfterStr != "" {
		createdAfter, err := time.Parse(time.RFC3339, createdAfterStr)
		if err == nil {
			params.CreatedAfter = &createdAfter
		}
	}

	if createdBeforeStr := c.Query("created_before"); createdBeforeStr != "" {
		createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr)
		if err == nil {
			params.CreatedBefore = &createdBefore
		}
	}

	walletAddress, err := controller.IWalletAddressService.GetList(params)
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

func (controller *WalletAddressController) AllocateWalletAddresses(c *gin.Context) {
	var req requests.AllocateWalletAddress
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

	walletAddress, err := controller.IWalletAddressService.AllocateWalletAddresses(&req)
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

	transactions, err := controller.IWalletAddressService.GetTransactions(walletAddress,
		uint(c.GetInt("page")),
		uint(c.GetInt("limit")),
	)

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

func (controller *WalletAddressController) GetMyAllocated(c *gin.Context) {
	auth := utils.GetAuthData(c)
	wallets, err := controller.IWalletAddressService.GetAllocatedByUser(auth.OwnerId)
	if err != nil {
		response.Api(c).SetMessage(err.Error()).Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"wallet_addresses": wallets,
		}).Send()
}

func (controller *WalletAddressController) GetOwnedTransactions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		response.Api(c).Send()
		return
	}

	walletAddress, err := controller.IWalletAddressService.GetByUuid(&id)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).Send()
		return
	}

	auth := utils.GetAuthData(c)
	if walletAddress.AllocatedToUserID == nil || *walletAddress.AllocatedToUserID != auth.OwnerId {
		response.Api(c).SetStatusCode(http.StatusForbidden).SetMessage(i18n.Localize(c.GetString("locale"), "request-unauthorized")).Send()
		return
	}

	transactions, err := controller.IWalletAddressService.GetTransactions(
		walletAddress,
		uint(c.GetInt("page")),
		uint(c.GetInt("limit")),
	)
	if err != nil {
		response.Api(c).SetMessage(err.Error()).Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"transactions": transactions,
		}).Send()
}

func (controller *WalletAddressController) Webhook(c *gin.Context) {
	hash := c.Param("hash")
	fmt.Println("received transaction", hash)
}
