package repositories

import (
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"fmt"

	"github.com/google/uuid"
)

type IDepositRepository interface {
	TransactionsExist(txHash string) (bool, error)
	Create(deposit *models.Deposits) error
	GetList(params *scopes.QueryBuilderModel) ([]*models.Deposits, int64, error)
	GetByUuid(id *uuid.UUID) (*models.Deposits, error)
	GetByPaymentRequest(paymentRequestID uint) ([]*models.Deposits, error)
	DashboardSeries(userID *uint, days int) ([]map[string]interface{}, error)
	DashboardByBlockchain(userID *uint) ([]map[string]interface{}, error)
	CountByStatus(userID *uint) (map[string]int64, error)
	SumAmount(userID *uint) (float64, error)
}

type DepositRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *DepositRepository) TransactionsExist(txHash string) (bool, error) {
	var count int64
	result := repository.IDatabaseHandler.GetClient().
		Model(&models.Deposits{}).
		Where("transaction_hash = ?", txHash).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check if transaction exists: %s", result.Error.Error())
	}
	return count > 0, nil
}

func (repository *DepositRepository) Create(deposit *models.Deposits) error {
	result := repository.IDatabaseHandler.GetClient().Create(deposit)
	if result.Error != nil {
		return fmt.Errorf("failed to insert transaction into deposits table: %s", result.Error.Error())
	}
	return nil
}

func (repository *DepositRepository) GetList(params *scopes.QueryBuilderModel) ([]*models.Deposits, int64, error) {
	var items []*models.Deposits
	var count int64
	query := repository.IDatabaseHandler.GetClient().Model(&models.Deposits{}).
		Preload("Blockchain").
		Preload("WalletAddress").
		Preload("PaymentRequest")
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
	if err := query.Order("created_at desc").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (repository *DepositRepository) GetByUuid(id *uuid.UUID) (*models.Deposits, error) {
	var item models.Deposits
	err := repository.IDatabaseHandler.GetClient().
		Preload("Blockchain").
		Preload("WalletAddress").
		Preload("PaymentRequest").
		First(&item, "uuid = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repository *DepositRepository) GetByPaymentRequest(paymentRequestID uint) ([]*models.Deposits, error) {
	var items []*models.Deposits
	err := repository.IDatabaseHandler.GetClient().
		Preload("Blockchain").
		Where("payment_request_id = ?", paymentRequestID).
		Order("created_at desc").
		Find(&items).Error
	return items, err
}

func (repository *DepositRepository) DashboardSeries(userID *uint, days int) ([]map[string]interface{}, error) {
	type row struct {
		Date   string  `json:"date"`
		Count  int64   `json:"count"`
		Amount float64 `json:"amount"`
	}
	var rows []row
	query := repository.IDatabaseHandler.GetClient().Model(&models.Deposits{}).
		Select("to_char(created_at, 'YYYY-MM-DD') as date, count(*) as count, coalesce(sum(amount), 0) as amount").
		Where("created_at >= now() - interval '1 day' * ?", days).
		Group("to_char(created_at, 'YYYY-MM-DD')").
		Order("date asc")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(rows))
	for _, item := range rows {
		result = append(result, map[string]interface{}{
			"date":   item.Date,
			"count":  item.Count,
			"amount": item.Amount,
		})
	}
	return result, nil
}

func (repository *DepositRepository) DashboardByBlockchain(userID *uint) ([]map[string]interface{}, error) {
	type row struct {
		Blockchain string  `json:"blockchain"`
		Count      int64   `json:"count"`
		Amount     float64 `json:"amount"`
	}
	var rows []row
	query := repository.IDatabaseHandler.GetClient().Table("deposits").
		Select("coalesce(blockchains.name, 'unknown') as blockchain, count(deposits.id) as count, coalesce(sum(deposits.amount), 0) as amount").
		Joins("left join blockchains on blockchains.id = deposits.blockchain_id").
		Where("deposits.deleted_at is null").
		Group("blockchains.name")
	if userID != nil {
		query = query.Where("deposits.user_id = ?", *userID)
	}
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(rows))
	for _, item := range rows {
		result = append(result, map[string]interface{}{
			"blockchain": item.Blockchain,
			"count":      item.Count,
			"amount":     item.Amount,
		})
	}
	return result, nil
}

func (repository *DepositRepository) CountByStatus(userID *uint) (map[string]int64, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	query := repository.IDatabaseHandler.GetClient().Model(&models.PaymentRequest{}).
		Select("status, count(*) as count").
		Group("status")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := map[string]int64{}
	for _, item := range rows {
		result[item.Status] = item.Count
	}
	return result, nil
}

func (repository *DepositRepository) SumAmount(userID *uint) (float64, error) {
	var sum float64
	query := repository.IDatabaseHandler.GetClient().Model(&models.Deposits{}).Select("coalesce(sum(amount), 0)")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if err := query.Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}
