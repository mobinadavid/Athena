package repositories

import (
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type IIGPRepository interface {
	GetList(params *scopes.QueryBuilderModel) ([]*models.IGPModel, int64, error)
	GetByUuid(uuid *uuid.UUID) (*models.IGPModel, error)
	GetForVerification(uuid *uuid.UUID, issuerReferenceNumber string) (*models.IGPModel, error)
	Create(IGP *models.IGPModel) (*models.IGPModel, error)
	Update(IGP *models.IGPModel) (*models.IGPModel, error)
}

type IGPRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *IGPRepository) GetList(params *scopes.QueryBuilderModel) ([]*models.IGPModel, int64, error) {
	var igps []*models.IGPModel
	var count int64

	query := repository.IDatabaseHandler.GetClient().Model(&models.IGPModel{})

	// Apply filters using parameterized queries
	validFilters := utils.GetStructFieldNames(models.IGPModel{})
	namingStrategy := schema.NamingStrategy{}

	for key, value := range params.Filters {
		if validFilters[key] {
			query = query.Where(fmt.Sprintf("%s = ?", namingStrategy.ColumnName("", key)), value)
		}
	}

	// Apply created_at range filters
	if params.CreatedAfter != nil {
		query = query.Where("created_at >= ?", params.CreatedAfter)
	}
	if params.CreatedBefore != nil {
		query = query.Where("created_at <= ?", params.CreatedBefore)
	}

	// Apply sorting using safe methods
	if params.SortBy != "" {
		sortOrder := "asc"
		if params.SortOrder == "desc" {
			sortOrder = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", namingStrategy.ColumnName("", params.SortBy), sortOrder))
	}

	// Get total count before pagination
	query.Count(&count)

	// Apply pagination
	if params.Page != 0 && params.Limit != 0 {
		query = query.Scopes(scopes.PaginateScope(params.Page, params.Limit))
	}

	// Execute the query
	result := query.Find(&igps)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, fmt.Errorf("igps get list failed: %s", result.Error.Error())
	}

	return igps, count, nil
}

func (repository *IGPRepository) GetByUuid(uuid *uuid.UUID) (*models.IGPModel, error) {
	var igp models.IGPModel
	result := repository.IDatabaseHandler.GetClient().Model(models.IGPModel{}).First(&igp, "uuid = ?", uuid)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("igp get by uuid failed: %s", result.Error.Error())
	}

	return &igp, nil
}

func (repository *IGPRepository) GetForVerification(uuid *uuid.UUID, issuerReferenceNumber string) (*models.IGPModel, error) {
	var igp models.IGPModel
	result := repository.IDatabaseHandler.GetClient().Model(models.IGPModel{}).
		Where("uuid = ?", uuid).
		Where("issuer_reference_number != ? OR issuer_reference_number IS NULL", issuerReferenceNumber).
		Where("status = ?", "pending").
		First(&igp)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("igp get by uuid failed: %s", result.Error.Error())
	}

	return &igp, nil
}

func (repository *IGPRepository) Create(IGP *models.IGPModel) (*models.IGPModel, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&IGP)
	if result.Error != nil {
		return nil, fmt.Errorf("igp creation failed: %s", result.Error.Error())
	}

	return IGP, nil
}

func (repository *IGPRepository) Update(IGP *models.IGPModel) (*models.IGPModel, error) {
	result := repository.IDatabaseHandler.GetClient().Save(&IGP)
	if result.Error != nil {
		return nil, fmt.Errorf("igp update failed: %s", result.Error.Error())
	}
	return IGP, nil
}
