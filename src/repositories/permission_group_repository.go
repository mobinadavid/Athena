package repositories

import (
	"athena/src/api/errs"
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type IPermissionGroupRepository interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(uuid *uuid.UUID, includeRelation ...string) (*models.PermissionGroupModel, error)
	Create(permission *models.PermissionGroupModel) (*models.PermissionGroupModel, error)
	Update(group *models.PermissionGroupModel) (*models.PermissionGroupModel, error)
	Delete(group *models.PermissionGroupModel) error
}

type PermissionGroupRepository struct {
	DatabaseHandler *database.Database
}

func (repository *PermissionGroupRepository) Create(permission *models.PermissionGroupModel) (*models.PermissionGroupModel, error) {
	result := repository.DatabaseHandler.GetClient().Create(&permission)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "idx_permission_groups_name" {
				return nil, errs.DuplicatePermissionGroupName
			}
		}
		return nil, fmt.Errorf("permission group creation failed: %s", result.Error.Error())
	}
	return permission, nil
}

func (repository *PermissionGroupRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.PermissionGroupModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Preload("Permissions").Model(results)

	// Apply pagination, filtering, and sorting using the BuilderModel
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("permission group list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("permission group list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

func (repository *PermissionGroupRepository) GetByUuid(uuid *uuid.UUID, includeRelations ...string) (*models.PermissionGroupModel, error) {
	var permission models.PermissionGroupModel
	result := repository.DatabaseHandler.GetClient()
	for _, relation := range includeRelations {
		if relation != "" {
			result = result.Preload(relation)
		}
	}
	result = result.First(&permission, "uuid = ?", uuid)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &permission, nil
}

func (repository *PermissionGroupRepository) Update(group *models.PermissionGroupModel) (*models.PermissionGroupModel, error) {
	err := repository.DatabaseHandler.GetClient().Select("*").
		Updates(group).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "idx_permission_groups_name" {
				return nil, errs.DuplicatePermissionGroupName
			}
		}
		return nil, err
	}
	err = repository.DatabaseHandler.GetClient().Model(group).Association("Permissions").Replace(group.Permissions)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (repository *PermissionGroupRepository) Delete(group *models.PermissionGroupModel) error {
	result := repository.DatabaseHandler.GetClient().Delete(&group)
	if result.Error != nil {
		return fmt.Errorf("group delete failed: %s", result.Error.Error())
	}
	return nil
}
