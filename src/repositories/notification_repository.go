package repositories

import (
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"

	"github.com/google/uuid"
)

type INotificationRepository interface {
	Create(notification *models.Notification) (*models.Notification, error)
	GetList(params *scopes.QueryBuilderModel) ([]*models.Notification, int64, error)
	GetByUuid(id *uuid.UUID) (*models.Notification, error)
	Save(notification *models.Notification) (*models.Notification, error)
	MarkRead(notification *models.Notification) error
	MarkAllRead(userID uint) error
	UnreadCount(userID uint) (int64, error)
}

type NotificationRepository struct {
	DatabaseHandler *database.Database
}

func (repository *NotificationRepository) Create(notification *models.Notification) (*models.Notification, error) {
	if err := repository.DatabaseHandler.GetClient().Create(notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

func (repository *NotificationRepository) GetList(params *scopes.QueryBuilderModel) ([]*models.Notification, int64, error) {
	var items []*models.Notification
	var count int64
	query := repository.DatabaseHandler.GetClient().Model(&models.Notification{})
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if params.Page != 0 && params.Limit != 0 {
		query = query.Scopes(scopes.PaginateScope(params.Page, params.Limit))
	}
	if err := query.Order("created_at desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (repository *NotificationRepository) GetByUuid(id *uuid.UUID) (*models.Notification, error) {
	var item models.Notification
	if err := repository.DatabaseHandler.GetClient().First(&item, "uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (repository *NotificationRepository) Save(notification *models.Notification) (*models.Notification, error) {
	if err := repository.DatabaseHandler.GetClient().Save(notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

func (repository *NotificationRepository) MarkRead(notification *models.Notification) error {
	return repository.DatabaseHandler.GetClient().Model(notification).Update("is_read", true).Error
}

func (repository *NotificationRepository) MarkAllRead(userID uint) error {
	return repository.DatabaseHandler.GetClient().Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

func (repository *NotificationRepository) UnreadCount(userID uint) (int64, error) {
	var count int64
	err := repository.DatabaseHandler.GetClient().Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
