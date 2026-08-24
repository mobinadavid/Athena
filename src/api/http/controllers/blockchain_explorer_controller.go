package controllers

import (
	"athena/src/api/http/requests"
	"athena/src/api/http/response"
	"athena/src/database/scopes"
	"athena/src/pkg/i18n"
	"athena/src/pkg/validator"
	"athena/src/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlockchainExplorerController struct {
	IBlockchainExplorerService services.IBlockchainExplorerService
}

func (controller *BlockchainExplorerController) GetList(c *gin.Context) {
	var blockchainExplorer *scopes.PaginatedModel
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

	blockchainExplorer, err := controller.IBlockchainExplorerService.GetList(params)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"blockchain_explorers": blockchainExplorer,
		}).Send()
	return
}

func (controller *BlockchainExplorerController) GetByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	//  to find the blockchainExplorer by UUID
	blockchainExplorer, err := controller.IBlockchainExplorerService.GetByUuid(&id)

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
			"blockchain_explorer:": blockchainExplorer,
		}).
		Send()
	return
}

func (controller *BlockchainExplorerController) Create(c *gin.Context) {
	var req requests.CreateBlockchainExplorerRequest

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

	blockchainExplorer, err := controller.IBlockchainExplorerService.Create(&req)

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
			"blockchain_explorer": blockchainExplorer,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()
}

func (controller *BlockchainExplorerController) Delete(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}
	err = controller.IBlockchainExplorerService.Delete(&id)

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

func (controller *BlockchainExplorerController) Update(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")

	// Parse the string to a UUID
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).Send()
		return
	}

	var req requests.CreateBlockchainExplorerRequest

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

	blockchainExplorer, err := controller.IBlockchainExplorerService.Update(&id, &req)

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
			"blockchain_explorer": blockchainExplorer,
		}).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()

}
