package services

import (
	"athena/src/api/errs"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/repositories"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type INotificationService interface {
	Create(userID uint, notificationType, title, body string, data map[string]interface{}) (*models.Notification, error)
	GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error)
	MarkRead(userID uint, id *uuid.UUID) error
	MarkAllRead(userID uint) error
	UnreadCount(userID uint) (int64, error)
}

type NotificationService struct {
	NotificationRepository repositories.INotificationRepository
}

func (service *NotificationService) Create(userID uint, notificationType, title, body string, data map[string]interface{}) (*models.Notification, error) {
	var payload datatypes.JSON
	if data != nil {
		encoded, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		payload = encoded
	}

	return service.NotificationRepository.Create(&models.Notification{
		UserID: userID,
		Type:   notificationType,
		Title:  title,
		Body:   body,
		Data:   payload,
	})
}

func (service *NotificationService) GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error) {
	items, count, err := service.NotificationRepository.GetList(params)
	if err != nil {
		return nil, err
	}
	return paginatedResult(params, items, count), nil
}

func (service *NotificationService) MarkRead(userID uint, id *uuid.UUID) error {
	notification, err := service.NotificationRepository.GetByUuid(id)
	if err != nil {
		return errs.RecordNotFound
	}
	if notification.UserID != userID {
		return errs.ErrForbiddenResource
	}
	return service.NotificationRepository.MarkRead(notification)
}

func (service *NotificationService) MarkAllRead(userID uint) error {
	return service.NotificationRepository.MarkAllRead(userID)
}

func (service *NotificationService) UnreadCount(userID uint) (int64, error) {
	return service.NotificationRepository.UnreadCount(userID)
}
