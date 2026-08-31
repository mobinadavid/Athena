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

// IAdminRepository interface defines the methods to interact with the admin data store.
type IAdminRepository interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetExcelData(ids []string, filters map[string]any) ([]*models.AdminModel, error)
	GetExcelExport(builder *scopes.BuilderModel, ids []string) ([]*models.AdminModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.AdminModel, error)
	GetById(adminID uint) (*models.AdminModel, error)
	GetByMobile(mobile string) (*models.AdminModel, error)
	GetByNationalIdentityCode(nationalIdentityCode string) (*models.AdminModel, error)
	GetByUsername(username string) (*models.AdminModel, error)
	Create(admin *models.AdminModel) (*models.AdminModel, error)
	Save(admin *models.AdminModel) (*models.AdminModel, error)
	UpdateWithRole(admin *models.AdminModel) (*models.AdminModel, error)
	Delete(admin *models.AdminModel) error
	Update(admin *models.AdminModel) (*models.AdminModel, error)
}

// AdminRepository struct implements the AdminRepository interface.
type AdminRepository struct {
	DatabaseHandler *database.Database
}

// GetList retrieve all admin
func (repository *AdminRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.AdminModel
	builder.Relations = append(builder.Relations, "Roles")
	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)
	var err error
	nameStrs := make([]string, 0)
	if _, exist := builder.M2MFilters["Roles.name"]; exist {
		names := builder.M2MFilters["Roles.name"]
		for _, name := range names {
			nameStrs = append(nameStrs, name.(string))
		}
	}

	db, err = builder.QueryBuilderScopeWithRoleAdminM2M(db, nameStrs)
	if err != nil {
		return nil, err
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("admin list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

func (repository *AdminRepository) GetExcelData(ids []string, filters map[string]any) ([]*models.AdminModel, error) {
	var admins []*models.AdminModel
	query := repository.DatabaseHandler.GetClient().Order("first_name ASC")

	if len(ids) > 0 {
		query = query.Where("uuid IN (?)", ids)
	}

	for key, value := range filters {
		query = query.Where(key, value)
	}

	err := query.Find(&admins).Error
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (repository *AdminRepository) GetExcelExport(builder *scopes.BuilderModel, ids []string) ([]*models.AdminModel, error) {
	var results []*models.AdminModel
	builder.Relations = append(builder.Relations, "Roles")
	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)
	var err error
	if len(ids) > 0 {
		db = db.Where("uuid IN (?)", ids)
	}
	nameStrs := make([]string, 0)
	if _, exist := builder.M2MFilters["Roles.name"]; exist {
		names := builder.M2MFilters["Roles.name"]
		for _, name := range names {
			nameStrs = append(nameStrs, name.(string))
		}
	}
	db, err = builder.QueryBuilderScopeWithRoleAdminM2M(db, nameStrs)
	if err != nil {
		return nil, fmt.Errorf("admin list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	err = db.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("admin get all retrieval failed: %s", err.Error())
	}

	if len(results) == 0 {
		return nil, errs.RecordNotFound
	}

	return results, nil
}

// GetById retrieve an admin by id
func (repository *AdminRepository) GetById(adminID uint) (*models.AdminModel, error) {
	var admin models.AdminModel
	result := repository.DatabaseHandler.GetClient().Preload("Roles.PermissionGroups.Permissions").First(&admin, "id = ?", adminID)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &admin, nil
}

// GetByUuid gets an admin by uuid.
func (repository *AdminRepository) GetByUuid(uuid *uuid.UUID) (*models.AdminModel, error) {
	var admin models.AdminModel
	result := repository.DatabaseHandler.GetClient().Preload("Roles.PermissionGroups.Permissions").First(&admin, "uuid = ?", uuid)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &admin, nil
}

// GetByMobile gets an admin by mobile.
func (repository *AdminRepository) GetByMobile(mobile string) (*models.AdminModel, error) {
	var admin models.AdminModel
	result := repository.DatabaseHandler.GetClient().First(&admin, "mobile = ?", mobile)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &admin, nil
}

// GetByNationalIdentityCode retrieve an admin by national-identity-code.
func (repository *AdminRepository) GetByNationalIdentityCode(nationalIdentityCode string) (*models.AdminModel, error) {
	var admin models.AdminModel
	result := repository.DatabaseHandler.GetClient().First(&admin, "national_identity_code = ?", nationalIdentityCode)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &admin, nil
}

func (repository *AdminRepository) GetByUsername(username string) (*models.AdminModel, error) {
	var admin models.AdminModel
	result := repository.DatabaseHandler.GetClient().First(&admin, "username = ?", username)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &admin, nil
}

// Create inserts a new admin into the database
func (repository *AdminRepository) Create(admin *models.AdminModel) (*models.AdminModel, error) {
	result := repository.DatabaseHandler.GetClient().Create(&admin)
	if result.Error != nil {
		// Check for Postgres duplicate key error
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == errs.DuplicateKey {
				return nil, errs.DuplicateUsername
			}
		}
		return nil, fmt.Errorf("admin creation failed: %s", result.Error.Error())
	}
	return admin, nil
}

// Save update an admin
func (repository *AdminRepository) Save(admin *models.AdminModel) (*models.AdminModel, error) {
	result := repository.DatabaseHandler.GetClient().Save(&admin)
	if result.Error != nil {
		return nil, fmt.Errorf("admin save failed: %s", result.Error.Error())
	}
	return admin, nil
}

// UpdateWithRole update an admin
func (repository *AdminRepository) UpdateWithRole(admin *models.AdminModel) (*models.AdminModel, error) {
	// Start a transaction
	tx := database.GetInstance().GetClient().Begin()

	// Save the admin
	if err := tx.Save(admin).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Save roles
	if err := tx.Model(admin).Association("Roles").Replace(admin.Roles); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Delete delete an admin
func (repository *AdminRepository) Delete(admin *models.AdminModel) error {
	result := repository.DatabaseHandler.GetClient().Delete(&admin)
	if result.Error != nil {
		return fmt.Errorf("admin delete failed: %s", result.Error.Error())
	}
	return nil
}

// Update update a admin changed field
func (repository *AdminRepository) Update(admin *models.AdminModel) (*models.AdminModel, error) {
	result := repository.DatabaseHandler.GetClient().Model(admin).Updates(admin)
	if result.Error != nil {
		return nil, fmt.Errorf("admin update failed: %w", result.Error)
	}
	return admin, nil
}
