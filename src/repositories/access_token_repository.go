package repositories

import (
	"athena/src/api/errs"
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/hash"
	"athena/src/models"
	"athena/src/pkg/utils"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type IAccessTokenRepository interface {
	GetAll(ownerID uint, ownerType string) ([]*models.AccessTokenModel, error)
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetActiveTokens(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(accessTokenUuid *uuid.UUID) (*models.AccessTokenModel, error)
	Create(accessToken *models.AccessTokenModel) (*models.AccessTokenModel, error)
	UpdateLastUsedAt(accessToken *models.AccessTokenModel, timestamp time.Time) (*models.AccessTokenModel, error)
	RefreshAccessTokens(accessToken *models.AccessTokenModel, uuid uuid.UUID, newAccessToken, newRefreshToken []byte, accessTokenExpiresAt, refreshTokenExpiresAt time.Time) (*models.AccessTokenModel, error)
	Delete(accessToken *models.AccessTokenModel) error
	DeleteMany(accessTokens []*models.AccessTokenModel) error
}

type AccessTokenRepository struct {
	DatabaseHandler *database.Database
}

func (repository *AccessTokenRepository) GetAll(ownerID uint, ownerType string) ([]*models.AccessTokenModel, error) {
	var results []*models.AccessTokenModel
	res := repository.DatabaseHandler.GetClient().Where("owner_id = ?", ownerID).Where("owner_type = ?", ownerType).Find(&results)
	if res.Error != nil {
		return nil, fmt.Errorf("access token list retrieval failed: %s", res.Error)
	}
	return results, nil
}

func (repository *AccessTokenRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.AccessTokenModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)

	// Apply pagination, filtering, and sorting using the BuilderModel
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("access token list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("access token list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

func (repository *AccessTokenRepository) GetActiveTokens(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.AccessTokenModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)

	// get active token only
	db = db.Where("access_token_expires_at > ?", time.Now())

	// Apply pagination, filtering, and sorting using the BuilderModel
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("access token list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("access token list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

func (repository *AccessTokenRepository) Create(accessToken *models.AccessTokenModel) (*models.AccessTokenModel, error) {
	var err error

	accessToken.AccessToken, err = hash.GetInstance().Generate(accessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	accessToken.RefreshToken, err = hash.GetInstance().Generate(accessToken.RefreshToken)
	if err != nil {
		return nil, err
	}

	result := repository.DatabaseHandler.GetClient().Create(&accessToken)
	if result.Error != nil {
		return nil, fmt.Errorf("access token creation failed: %s", result.Error.Error())
	}
	return accessToken, nil
}

func (repository *AccessTokenRepository) GetByUuid(accessTokenUuid *uuid.UUID) (*models.AccessTokenModel, error) {
	var accessToken models.AccessTokenModel

	query := repository.DatabaseHandler.GetClient()

	result := query.First(&accessToken, "uuid = ?", accessTokenUuid)

	if result.Error != nil {
		if utils.CheckError(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFound
		}
		return nil, result.Error
	}

	return &accessToken, nil
}

func (repository *AccessTokenRepository) UpdateLastUsedAt(accessToken *models.AccessTokenModel, timestamp time.Time) (*models.AccessTokenModel, error) {
	result := repository.DatabaseHandler.GetClient().Model(accessToken).Update("LastUsedAt", timestamp)
	if result.Error != nil {
		return nil, result.Error
	}
	accessToken.LastUsedAt = &timestamp
	return accessToken, nil
}

func (repository *AccessTokenRepository) RefreshAccessTokens(accessToken *models.AccessTokenModel, uuid uuid.UUID, newAccessToken, newRefreshToken []byte, accessTokenExpiresAt, refreshTokenExpiresAt time.Time) (*models.AccessTokenModel, error) {
	var err error

	newAccessTokenHash, err := hash.GetInstance().Generate(newAccessToken)
	if err != nil {
		return nil, err
	}

	newRefreshTokenHash, err := hash.GetInstance().Generate(newRefreshToken)
	if err != nil {
		return nil, err
	}

	result := repository.DatabaseHandler.GetClient().Model(&accessToken).Updates(models.AccessTokenModel{
		Uuid:                  uuid,
		AccessToken:           newAccessTokenHash,
		RefreshToken:          newRefreshTokenHash,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return accessToken, nil
}

func (repository *AccessTokenRepository) Delete(accessToken *models.AccessTokenModel) error {
	err := repository.DatabaseHandler.GetClient().Delete(&accessToken)
	if err.Error != nil {
		return fmt.Errorf("access tokens destroy failed: %s", err.Error.Error())
	}
	return nil
}

func (repository *AccessTokenRepository) DeleteMany(accessTokens []*models.AccessTokenModel) error {
	err := repository.DatabaseHandler.GetClient().Delete(&accessTokens)
	if err.Error != nil {
		return fmt.Errorf("access tokens destroy failed: %s", err.Error.Error())
	}
	return nil
}
