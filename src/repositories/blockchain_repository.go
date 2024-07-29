package repositories

import (
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBlockchainRepository interface {
	GetList() ([]*models.Blockchain, error)
	GetByUuid(uuid *uuid.UUID) (*models.Blockchain, error)
	GetByName(name string) (*models.Blockchain, error)
	Create(blockchain *models.Blockchain) (*models.Blockchain, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *models.Blockchain) (*models.Blockchain, error)
}

type BlockchainRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *BlockchainRepository) GetList() ([]*models.Blockchain, error) {
	var blockchain []*models.Blockchain

	result := repository.IDatabaseHandler.GetClient().Preload("WalletAddresses").Preload("BlockchainExplorers").Model(&models.Blockchain{}).Scopes()
	result = result.Find(&blockchain)

	if result.Error != nil {
		return nil, fmt.Errorf("blockchain get list failed: %s", result.Error.Error())
	}

	return blockchain, nil

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

func (repository *BlockchainRepository) GetCount() (int64, error) {
	var count int64

	result := repository.IDatabaseHandler.GetClient().Model(&models.Blockchain{})
	result = result.Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("blockchain get count failed: %s", result.Error.Error())
	}

	return count, nil
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
