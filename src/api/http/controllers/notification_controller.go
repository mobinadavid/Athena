package controllers

import (
	"athena/src/api/http/response"
	"athena/src/pkg/i18n"
	"athena/src/pkg/utils"
	"athena/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {
	NotificationService services.INotificationService
}

func (controller *NotificationController) GetList(c *gin.Context) {
	auth := utils.GetAuthData(c)
	params := listParams(c)
	params.UserID = &auth.OwnerId
	items, err := controller.NotificationService.GetList(params)
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	unread, _ := controller.NotificationService.UnreadCount(auth.OwnerId)
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		SetData(map[string]interface{}{
			"notifications": items,
			"unread_count":  unread,
		}).
		Send()
}

func (controller *NotificationController) MarkRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	auth := utils.GetAuthData(c)
	if err := controller.NotificationService.MarkRead(auth.OwnerId, &id); err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()
}

func (controller *NotificationController) MarkAllRead(c *gin.Context) {
	auth := utils.GetAuthData(c)
	if err := controller.NotificationService.MarkAllRead(auth.OwnerId); err != nil {
		response.Api(c).SetMessage(i18n.Localize(c.GetString("locale"), err.Error())).Send()
		return
	}
	response.Api(c).
		SetStatusCode(http.StatusOK).
		SetMessage(i18n.Localize(c.GetString("locale"), "request-successful")).
		Send()
}
