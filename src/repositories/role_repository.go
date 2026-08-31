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

type IRoleRepository interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetExcelData(ids []string, filters map[string]any) ([]*models.RoleModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.RoleModel, error)
	GetAssociateAdmin(role *models.RoleModel) ([]*models.AdminModel, error)
	GetAssociateUser(role *models.RoleModel) ([]*models.UserModel, error)
	Create(role *models.RoleModel) (*models.RoleModel, error)
	Update(role *models.RoleModel) (*models.RoleModel, error)
	DeletePermissions(role *models.RoleModel) error
	Delete(role *models.RoleModel) error
	GetByName(name string) (*models.RoleModel, error)
}

type RoleRepository struct {
	DatabaseHandler *database.Database
}

func (repository *RoleRepository) Create(role *models.RoleModel) (*models.RoleModel, error) {
	result := repository.DatabaseHandler.GetClient().Create(&role)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "idx_roles_name" {
				return nil, errs.DuplicateRoleName
			}
		}
		return nil, fmt.Errorf("role creation failed: %s", result.Error.Error())
	}
	return role, nil
}

func (repository *RoleRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.RoleModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results).Preload("PermissionGroups").Preload("Admins")

	// Apply pagination, filtering, and sorting using the BuilderModel
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("role list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("role list retrieval failed: %s", err.Error())
	}

	for _, role := range results {
		role.TotalAssignees = int64(len(role.Admins))
		role.Admins = nil
	}

	return paginateModel, nil
}

func (repository *RoleRepository) GetExcelData(ids []string, filters map[string]any) ([]*models.RoleModel, error) {
	var results []*models.RoleModel

	query := repository.DatabaseHandler.GetClient().Order("title ASC")

	if len(ids) > 0 {
		query = query.Where("uuid IN (?)", ids)
	}

	if len(filters) > 0 {
		for filter, value := range filters {
			query = query.Where(filter, value)
		}
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("role get all retrieval failed: %s", err.Error())
	}

	return results, nil
}

func (repository *RoleRepository) GetByUuid(uuid *uuid.UUID) (*models.RoleModel, error) {
	var role models.RoleModel
	result := repository.DatabaseHandler.GetClient().Preload("PermissionGroups").First(&role, "uuid = ?", uuid).Preload("Admins")
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	role.TotalAssignees = int64(len(role.Admins))
	role.Admins = nil

	return &role, nil
}

func (repository *RoleRepository) GetAssociateUser(role *models.RoleModel) ([]*models.UserModel, error) {
	var users []*models.UserModel
	if err := repository.DatabaseHandler.GetClient().Model(&role).Association("Users").Find(&users); err != nil {
		return nil, fmt.Errorf("user for role get failed: %s", err.Error())
	}
	return users, nil
}

func (repository *RoleRepository) GetAssociateAdmin(role *models.RoleModel) ([]*models.AdminModel, error) {
	var admins []*models.AdminModel
	if err := repository.DatabaseHandler.GetClient().Model(&role).Association("Admins").Find(&admins); err != nil {
		return nil, fmt.Errorf("admin for role get failed: %s", err.Error())
	}
	return admins, nil
}

func (repository *RoleRepository) Update(role *models.RoleModel) (*models.RoleModel, error) {
	result := repository.DatabaseHandler.GetClient().Select("*").Updates(&role)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "idx_roles_name" {
				return nil, errs.DuplicateRoleName
			}
		}
		return nil, fmt.Errorf("role update failed: %s", result.Error.Error())
	}
	err := repository.DatabaseHandler.GetClient().Model(role).Association("PermissionGroups").Replace(role.PermissionGroups)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (repository *RoleRepository) Delete(role *models.RoleModel) error {
	result := repository.DatabaseHandler.GetClient().Delete(&role)
	if result.Error != nil {
		return fmt.Errorf("role delete failed: %s", result.Error.Error())
	}
	return nil
}

func (repository *RoleRepository) DeletePermissions(role *models.RoleModel) error {
	if err := repository.DatabaseHandler.GetClient().Model(&role).Association("PermissionGroups").Clear(); err != nil {
		return fmt.Errorf("failed to clear permissions: %s", err.Error())
	}
	return nil
}

func (repository *RoleRepository) GetByName(name string) (*models.RoleModel, error) {
	var role models.RoleModel

	result := repository.DatabaseHandler.GetClient().Where("name = ?", name).First(&role)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, fmt.Errorf("failed to retrieve role by name '%s': %w", name, result.Error)
	}

	return &role, nil
}
