package controllers

import (
	"athena/src/api/http/requests"
	"athena/src/api/http/response"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/i18n"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentController struct {
	PaymentService services.IPaymentService
}

func (controller *PaymentController) Create(c *gin.Context) {
	var req requests.AllocateWalletAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Api(c).Send()
		return
	}
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetErrors(err).Send()
		return
	}

	auth := utils.GetAuthData(c)
	item, err := controller.PaymentService.Create(auth.OwnerId, &req)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}

	response.Api(c).
		SetStatusCode(http.StatusCreated).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"payment_request": item,
		}).Send()
}

func (controller *PaymentController) GetList(c *gin.Context) {
	params := listParams(c)
	auth := utils.GetAuthData(c)
	if auth.OwnerType == models.UserRole {
		params.UserID = &auth.OwnerId
	}
	if status := c.Query("status"); status != "" {
		params.Filters["status"] = status
	}

	items, err := controller.PaymentService.GetList(params)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"payment_requests": items}).
		Send()
}

func (controller *PaymentController) GetByUuid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	item, err := controller.PaymentService.GetByUuid(scopedUserID(c), &id)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"payment_request": item}).
		Send()
}

func (controller *PaymentController) GetTransactions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	items, err := controller.PaymentService.GetTransactions(scopedUserID(c), &id, uint(c.GetInt("page")), uint(c.GetInt("limit")))
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"transactions": items}).
		Send()
}

func listParams(c *gin.Context) *scopes.QueryBuilderModel {
	return &scopes.QueryBuilderModel{
		Page:      uint(c.GetInt("page")),
		Limit:     uint(c.GetInt("limit")),
		SortBy:    c.GetString("sort_by"),
		SortOrder: c.GetString("sort_order"),
		Filters:   map[string]interface{}{},
	}
}

func scopedUserID(c *gin.Context) *uint {
	auth := utils.GetAuthData(c)
	if auth.OwnerType == models.UserRole {
		return &auth.OwnerId
	}
	return nil
}
