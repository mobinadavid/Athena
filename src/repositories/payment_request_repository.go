package repositories

import (
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"fmt"

	"github.com/google/uuid"
)

type IPaymentRequestRepository interface {
	Create(request *models.PaymentRequest) (*models.PaymentRequest, error)
	GetByUuid(uuid *uuid.UUID, preloads ...string) (*models.PaymentRequest, error)
	GetByID(id uint, preloads ...string) (*models.PaymentRequest, error)
	GetList(params *scopes.QueryBuilderModel) ([]*models.PaymentRequest, int64, error)
	Update(request *models.PaymentRequest) (*models.PaymentRequest, error)
	GetExpiredPending() ([]*models.PaymentRequest, error)
	Save(request *models.PaymentRequest) (*models.PaymentRequest, error)
	Delete(id uint) error
}

type PaymentRequestRepository struct {
	DatabaseHandler *database.Database
}

func (repository *PaymentRequestRepository) Create(request *models.PaymentRequest) (*models.PaymentRequest, error) {
	if err := repository.DatabaseHandler.GetClient().Create(request).Error; err != nil {
		return nil, fmt.Errorf("payment request create failed: %w", err)
	}
	return request, nil
}

func (repository *PaymentRequestRepository) GetByUuid(id *uuid.UUID, preloads ...string) (*models.PaymentRequest, error) {
	var request models.PaymentRequest
	query := repository.DatabaseHandler.GetClient()
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if err := query.First(&request, "uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (repository *PaymentRequestRepository) GetByID(id uint, preloads ...string) (*models.PaymentRequest, error) {
	var request models.PaymentRequest
	query := repository.DatabaseHandler.GetClient()
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if err := query.First(&request, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (repository *PaymentRequestRepository) GetList(params *scopes.QueryBuilderModel) ([]*models.PaymentRequest, int64, error) {
	var items []*models.PaymentRequest
	var count int64
	query := repository.DatabaseHandler.GetClient().Model(&models.PaymentRequest{}).Preload("Blockchain").Preload("WalletAddresses").Preload("Deposits")
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if status, ok := params.Filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if params.Page != 0 && params.Limit != 0 {
		query = query.Scopes(scopes.PaginateScope(params.Page, params.Limit))
	}
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := params.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (repository *PaymentRequestRepository) Update(request *models.PaymentRequest) (*models.PaymentRequest, error) {
	if err := repository.DatabaseHandler.GetClient().Save(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (repository *PaymentRequestRepository) Save(request *models.PaymentRequest) (*models.PaymentRequest, error) {
	return repository.Update(request)
}

func (repository *PaymentRequestRepository) GetExpiredPending() ([]*models.PaymentRequest, error) {
	var items []*models.PaymentRequest
	err := repository.DatabaseHandler.GetClient().
		Preload("WalletAddresses").
		Where("status = ? AND expires_at IS NOT NULL AND expires_at < now()", models.PaymentStatusPending).
		Find(&items).Error
	return items, err
}

func (repository *PaymentRequestRepository) Delete(id uint) error {
	return repository.DatabaseHandler.GetClient().Delete(&models.PaymentRequest{}, id).Error
}
