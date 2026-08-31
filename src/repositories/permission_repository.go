package repositories

import (
	"athena/src/api/errs"
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/utils"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPermissionRepository interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.PermissionModel, error)
	Create(permission *models.PermissionModel) (*models.PermissionModel, error)
}

type PermissionRepository struct {
	DatabaseHandler *database.Database
}

func (repository *PermissionRepository) Create(permission *models.PermissionModel) (*models.PermissionModel, error) {
	result := repository.DatabaseHandler.GetClient().Create(&permission)
	if result.Error != nil {
		return nil, fmt.Errorf("permission creation failed: %s", result.Error.Error())
	}
	return permission, nil
}

func (repository *PermissionRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.PermissionModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)

	// Apply pagination, filtering, and sorting using the BuilderModel
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("permission list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("permission list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

func (repository *PermissionRepository) GetByUuid(uuid *uuid.UUID) (*models.PermissionModel, error) {
	var permission models.PermissionModel
	result := repository.DatabaseHandler.GetClient().First(&permission, "uuid = ?", uuid)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &permission, nil
}
