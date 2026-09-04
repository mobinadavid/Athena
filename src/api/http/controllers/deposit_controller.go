package controllers

import (
	"athena/src/api/http/response"
	"athena/src/pkg/i18n"
	"athena/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepositController struct {
	DepositService services.IDepositService
}

func (controller *DepositController) GetList(c *gin.Context) {
	params := listParams(c)
	if userID := scopedUserID(c); userID != nil {
		params.UserID = userID
	}
	if status := c.Query("status"); status != "" {
		params.Filters["status"] = status
	}
	items, err := controller.DepositService.GetList(params)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"deposits": items}).
		Send()
}

func (controller *DepositController) GetByUuid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	item, err := controller.DepositService.GetByUuid(scopedUserID(c), &id)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"deposit": item}).
		Send()
}
