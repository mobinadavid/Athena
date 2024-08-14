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

type IBlockchainRepository interface {
	GetList(params *scopes.QueryBuilderModel) ([]*models.Blockchain, int64, error)
	GetByUuid(uuid *uuid.UUID) (*models.Blockchain, error)
	GetByName(name string) (*models.Blockchain, error)
	Create(blockchain *models.Blockchain) (*models.Blockchain, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *models.Blockchain) (*models.Blockchain, error)
}

type BlockchainRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *BlockchainRepository) GetList(params *scopes.QueryBuilderModel) ([]*models.Blockchain, int64, error) {
	var blockchains []*models.Blockchain
	var count int64
	query := repository.IDatabaseHandler.GetClient().Preload("BlockchainExplorers").Model(&models.Blockchain{})

	validFilters := utils.GetStructFieldNames(models.Blockchain{})
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
	result := query.Find(&blockchains)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, fmt.Errorf("blockchains get list failed: %s", result.Error.Error())
	}

	return blockchains, count, nil
}

func (repository *BlockchainRepository) GetByUuid(uuid *uuid.UUID) (*models.Blockchain, error) {
	var blockchain models.Blockchain

	result := repository.IDatabaseHandler.GetClient().Preload("WalletAddresses").Preload("BlockchainExplorers").First(&blockchain, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("blockchain get by uuid failed: %s", result.Error.Error())
	}

	return &blockchain, nil
}

func (repository *BlockchainRepository) GetByName(name string) (*models.Blockchain, error) {
	var blockchain models.Blockchain

	result := repository.IDatabaseHandler.GetClient().Scopes(scopes.IsActive()).First(&blockchain, "name = ?", name)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("blockchain get by name failed: %s", result.Error.Error())
	}

	return &blockchain, nil
}

func (repository *BlockchainRepository) Create(blockchain *models.Blockchain) (*models.Blockchain, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&blockchain)
	if result.Error != nil {
		return nil, fmt.Errorf("blockchain creation failed: %s", result.Error.Error())
	}

	return blockchain, nil
}

func (repository *BlockchainRepository) Delete(uuid *uuid.UUID) error {
	var blockchain models.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&blockchain, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("blockchain get by uuid failed: %s", result.Error.Error())
	}

	if err := repository.IDatabaseHandler.GetClient().Delete(&blockchain).Error; err != nil {
		return fmt.Errorf("failed to delete blockchain: %s", err)
	}

	return nil
}

func (repository *BlockchainRepository) Update(uuid *uuid.UUID, req *models.Blockchain) (*models.Blockchain, error) {
	var existingBlockchain models.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&existingBlockchain, "uuid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blockchain with UUID %s not found", uuid)
		}

		return nil, fmt.Errorf("failed to retrieve blockchain with UUID %s: %s", uuid, result.Error)
	}

	if err := repository.IDatabaseHandler.GetClient().Model(&existingBlockchain).Updates(req).Error; err != nil {
		return nil, fmt.Errorf("failed to update blockchain: %s", err)
	}

	return &existingBlockchain, nil
}
