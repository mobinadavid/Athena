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

// IUserRepository interface defines the methods to interact with the User data store.
type IUserRepository interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.UserModel, error)
	GetById(userID uint, includeRelations ...string) (*models.UserModel, error)
	GetByMobile(mobile string) (*models.UserModel, error)
	GetByNationalIdentityCode(nationalIdentityCode string, includeRelations ...string) (*models.UserModel, error)
	GetByMobileAndNationalIdentityCode(mobile, nationalIdentityCode string) (*models.UserModel, error)
	Create(user *models.UserModel) (*models.UserModel, error)
	Save(user *models.UserModel) (*models.UserModel, error)
	Update(user *models.UserModel) (*models.UserModel, error)
	Delete(user *models.UserModel) error
	FirstOrCreate(mobile string, nationalIdentityCode string) (*models.UserModel, error)
	UpdateOrCreate(search *models.UserModel, assign *models.UserModel) (*models.UserModel, error)
}

// UserRepository struct implements the UserRepository interface.
type UserRepository struct {
	DatabaseHandler *database.Database
}

// FirstOrCreate retrieve a user and create it if not exist create it first
func (repository *UserRepository) FirstOrCreate(mobile string, nationalIdentityCode string) (*models.UserModel, error) {
	user := &models.UserModel{
		Mobile:               mobile,
		NationalIdentityCode: nationalIdentityCode,
	}

	result := repository.DatabaseHandler.GetClient().FirstOrCreate(&user, models.UserModel{Mobile: mobile, NationalIdentityCode: nationalIdentityCode})

	if result.Error != nil {
		return nil, fmt.Errorf("user first or create failed: %s", result.Error.Error())
	}
	return user, nil
}

// GetList retrieve all users
func (repository *UserRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.UserModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)

	// Apply pagination, filtering, and sorting using the BuilderModel
	//builder.Relations = append(builder.Relations, "Wallet")
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("user list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("user list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

// GetById retrieve a user by id
func (repository *UserRepository) GetById(userID uint, includeRelations ...string) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient()
	for _, relation := range includeRelations {
		result.Preload(relation)

	}
	result = result.First(&user, "id = ?", userID)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByUuid retrieve a user by uuid
func (repository *UserRepository) GetByUuid(uuid *uuid.UUID) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().
		First(&user, "uuid = ?", uuid)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByMobile retrieve a user by mobile.
func (repository *UserRepository) GetByMobile(mobile string) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().First(&user, "mobile = ?", mobile)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByNationalIdentityCode gets a user by national-identity-code.
func (repository *UserRepository) GetByNationalIdentityCode(nationalIdentityCode string, includeRelations ...string) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient()
	for _, relation := range includeRelations {
		result = result.Preload(relation)
	}
	result = result.First(&user, "national_identity_code = ?", nationalIdentityCode)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByMobileAndNationalIdentityCode gets a user by mobile and national-identity-code.
func (repository *UserRepository) GetByMobileAndNationalIdentityCode(mobile, nationalIdentityCode string) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().First(&user, "national_identity_code = ? AND mobile = ?", nationalIdentityCode, mobile)
	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// Create inserts a new User into the database
func (repository *UserRepository) Create(user *models.UserModel) (*models.UserModel, error) {
	result := repository.DatabaseHandler.GetClient().Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user creation failed: %s", result.Error.Error())
	}
	return user, nil
}

// Save update a user's all field
func (repository *UserRepository) Save(user *models.UserModel) (*models.UserModel, error) {
	err := repository.DatabaseHandler.GetClient().
		Save(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user save failed: %s", err.Error())
	}
	return user, nil
}

// Update update a user's changed field
func (repository *UserRepository) Update(user *models.UserModel) (*models.UserModel, error) {
	result := repository.DatabaseHandler.GetClient().Updates(user)
	if result.Error != nil {
		return nil, fmt.Errorf("user update failed: %w", result.Error)
	}
	return user, nil
}

// Delete delete a user
func (repository *UserRepository) Delete(user *models.UserModel) error {
	result := repository.DatabaseHandler.GetClient().Delete(&user)
	if result.Error != nil {
		return fmt.Errorf("user delete failed: %s", result.Error.Error())
	}
	return nil
}

func (repository *UserRepository) UpdateOrCreate(search *models.UserModel, attr *models.UserModel) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().Where(search).Assign(attr).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user update or create failed: %s", result.Error.Error())
	}
	return &user, nil
}
