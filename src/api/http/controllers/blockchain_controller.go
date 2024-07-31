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
	"time"
)

type BlockchainController struct {
	IBlockchainService services.IBlockChainService
}

func (controller *BlockchainController) GetList(c *gin.Context) {
	var blockchain *scopes.PaginateModel
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

	blockchain, err := controller.IBlockchainService.GetList(params)

	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"blockchains": blockchain,
		}).Send()
	return
}

func (controller *BlockchainController) GetByUuid(c *gin.Context) {
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	//  to find the blockchain by UUID
	blockchain, err := controller.IBlockchainService.GetByUuid(&id)

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
			"blockchain:": blockchain,
		}).
		Send()
	return
}

func (controller *BlockchainController) Create(c *gin.Context) {
	var req requests.CreateBlockchainRequest

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

	blockchain, err := controller.IBlockchainService.Create(&req)

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
			"blockchain": blockchain,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()
}

func (controller *BlockchainController) Delete(c *gin.Context) {
	//Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}
	err = controller.IBlockchainService.Delete(&id)

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

func (controller *BlockchainController) Update(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	var req requests.CreateBlockchainRequest

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

	blockchain, err := controller.IBlockchainService.Update(&id, &req)

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
			"blockchain": blockchain,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}
