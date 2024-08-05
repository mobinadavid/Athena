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
	"time"
)

type IWalletAddressRepository interface {
	GetList(params *scopes.QueryBuilderModel) ([]*models.WalletAddress, int64, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Create(blockchain *models.WalletAddress) (*models.WalletAddress, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *models.WalletAddress) (*models.WalletAddress, error)
	GetAllocatedList() ([]*models.WalletAddress, error)
	GetUnallocatedWalletAddress(count int, blockchainId uint) ([]*models.WalletAddress, error)
	UpdateWalletAddressToAllocated(walletAddresses []*models.WalletAddress) error
}

type WalletAddressRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *WalletAddressRepository) GetList(params *scopes.QueryBuilderModel) ([]*models.WalletAddress, int64, error) {
	var walletAddress []*models.WalletAddress
	var count int64

	query := repository.IDatabaseHandler.GetClient().Preload("Blockchain").Model(&models.WalletAddress{})

	validFilters := utils.GetStructFieldNames(models.WalletAddress{})
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
	result := query.Find(&walletAddress)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, fmt.Errorf("walletAddresses get list failed: %s", result.Error.Error())
	}

	return walletAddress, count, nil

}

func (repository *WalletAddressRepository) GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error) {
	var walletAddress models.WalletAddress

	result := repository.IDatabaseHandler.GetClient().Preload("Blockchain").First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("walletAddress get by uuid failed: %s", result.Error.Error())
	}

	return &walletAddress, nil
}

func (repository *WalletAddressRepository) Create(walletAddress *models.WalletAddress) (*models.WalletAddress, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&walletAddress)
	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress creation failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) Delete(uuid *uuid.UUID) error {
	var walletAddress models.WalletAddress

	result := repository.IDatabaseHandler.GetClient().Preload("Blockchain").First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("walletAddress get by uuid failed: %s", result.Error.Error())
	}

	if err := repository.IDatabaseHandler.GetClient().Delete(&walletAddress).Error; err != nil {
		return fmt.Errorf("failed to delete walletAddress: %s", err)
	}

	return nil

}

func (repository *WalletAddressRepository) Update(uuid *uuid.UUID, req *models.WalletAddress) (*models.WalletAddress, error) {
	var walletAddress models.WalletAddress

	result := repository.IDatabaseHandler.GetClient().Preload("Blockchain").First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("walletAddress with UUID %s not found", uuid)
		}

		return nil, fmt.Errorf("failed to retrieve walletAddress with UUID %s: %s", uuid, result.Error)
	}

	if err := repository.IDatabaseHandler.GetClient().Session(&gorm.Session{FullSaveAssociations: true}).Model(&walletAddress).Updates(req).Error; err != nil {
		return nil, fmt.Errorf("failed to update walletAddress: %s", err)
	}

	return &walletAddress, nil
}

func (repository *WalletAddressRepository) UpdateWalletAddressToAllocated(walletAddresses []*models.WalletAddress) error {
	for _, walletAddress := range walletAddresses {
		walletAddress.AllocatedAt = time.Now()
		if err := repository.IDatabaseHandler.GetClient().Save(walletAddress).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repository *WalletAddressRepository) GetUnallocatedWalletAddress(count int, blockchainId uint) ([]*models.WalletAddress, error) {
	var walletAddresses []*models.WalletAddress
	// Use IsAllocated to filter records
	result := repository.IDatabaseHandler.GetClient().Scopes(scopes.IsNotAllocated()).Preload("Blockchain").Where("blockchain_id = ?", blockchainId).Limit(count)
	if err := result.Find(&walletAddresses).Error; err != nil {
		return nil, err
	}

	return walletAddresses, nil
}

func (repository *WalletAddressRepository) GetAllocatedList() ([]*models.WalletAddress, error) {
	var walletAddress []*models.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Scopes(scopes.IsAllocated()).Preload("Blockchain").Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}
