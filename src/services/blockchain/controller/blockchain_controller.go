package controller

import (
	"athena/src/api/http/response"
	"athena/src/database/scopes"
	"athena/src/pkg/i18n"
	"athena/src/pkg/validator"
	"athena/src/services/blockchain/request"
	"athena/src/services/blockchain/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type BlockchainController struct {
	IBlockchainService service.IBlockChainService
}

func (controller *BlockchainController) GetList(c *gin.Context) {
	var blockchain *scopes.PaginateModel

	blockchain, err := controller.IBlockchainService.GetList()

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
	var req request.CreateBlockchainRequest

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

	var req request.CreateBlockchainRequest

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
