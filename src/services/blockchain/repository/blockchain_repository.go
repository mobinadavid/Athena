package repository

import (
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/services/blockchain/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBlockchainRepository interface {
	GetList(page uint, limit uint) ([]*model.Blockchain, error)
	GetByUuid(uuid *uuid.UUID) (*model.Blockchain, error)
	Create(blockchain *model.Blockchain) (*model.Blockchain, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *model.Blockchain) (*model.Blockchain, error)
}

type BlockchainRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *BlockchainRepository) GetList(page uint, limit uint) ([]*model.Blockchain, error) {
	var blockchain []*model.Blockchain

	result := repository.IDatabaseHandler.GetClient().Scopes(
		scopes.PaginateScope(page, limit),
	)
	result = result.Where("is_active = ?", true).Find(&blockchain)

	if result.Error != nil {
		return nil, fmt.Errorf("blockchain get list failed: %s", result.Error.Error())
	}

	return blockchain, nil
}

func (repository *BlockchainRepository) GetByUuid(uuid *uuid.UUID) (*model.Blockchain, error) {
	var blockchain model.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&blockchain, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("blockchain get by uuid failed: %s", result.Error.Error())
	}

	return &blockchain, nil
}

func (repository *BlockchainRepository) Create(blockchain *model.Blockchain) (*model.Blockchain, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&blockchain)
	if result.Error != nil {
		return nil, fmt.Errorf("blockchain creation failed: %s", result.Error.Error())
	}

	return blockchain, nil
}

func (repository *BlockchainRepository) GetCount() (int64, error) {
	var count int64

	result := repository.IDatabaseHandler.GetClient().Model(&model.Blockchain{})
	result = result.Where("is_active = ?", true).Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("blockchain get count failed: %s", result.Error.Error())
	}

	return count, nil
}

func (repository *BlockchainRepository) Delete(uuid *uuid.UUID) error {
	var blockchain model.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&blockchain, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("blockchain get by uuid failed: %s", result.Error.Error())
	}

	if err := repository.IDatabaseHandler.GetClient().Delete(&blockchain).Error; err != nil {
		return fmt.Errorf("failed to delete blockchain: %s", err)
	}

	return nil
}

func (repository *BlockchainRepository) Update(uuid *uuid.UUID, req *model.Blockchain) (*model.Blockchain, error) {
	var existingBlockchain model.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&existingBlockchain, "uuid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category with UUID %s not found", uuid)
		}

		return nil, fmt.Errorf("failed to retrieve category with UUID %s: %s", uuid, result.Error)
	}

	if err := repository.IDatabaseHandler.GetClient().Model(&existingBlockchain).Updates(req).Error; err != nil {
		return nil, fmt.Errorf("failed to update FAQ: %s", err)
	}

	return &existingBlockchain, nil
}
