package controllers

import (
	"athena/src/api/http/response"
	"athena/src/models"
	"athena/src/pkg/i18n"
	"athena/src/pkg/utils"
	"athena/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	DashboardService services.IDashboardService
}

func (controller *DashboardController) Summary(c *gin.Context) {
	auth := utils.GetAuthData(c)
	var (
		data map[string]interface{}
		err  error
	)
	if auth.OwnerType == models.AdminRole {
		data, err = controller.DashboardService.AdminSummary()
	} else {
		data, err = controller.DashboardService.UserSummary(auth.OwnerId)
	}
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(data).
		Send()
}

func (controller *DashboardController) DepositsOverTime(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))
	series, err := controller.DashboardService.DepositsOverTime(scopedUserID(c), days)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"deposits_over_time": series}).
		Send()
}

func (controller *DashboardController) DepositsByBlockchain(c *gin.Context) {
	items, err := controller.DashboardService.DepositsByBlockchain(scopedUserID(c))
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{"deposits_by_blockchain": items}).
		Send()
}
