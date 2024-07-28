package repository

import (
	"athena/src/database"
	"athena/src/services/blockchain/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBlockchainRepository interface {
	GetList() ([]*model.Blockchain, error)
	GetByUuid(uuid *uuid.UUID) (*model.Blockchain, error)
	GetByName(name string) (*model.Blockchain, error)
	GetById(uint uint) (*model.Blockchain, error)
	Create(blockchain *model.Blockchain) (*model.Blockchain, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *model.Blockchain) (*model.Blockchain, error)
}

type BlockchainRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *BlockchainRepository) GetList() ([]*model.Blockchain, error) {
	var blockchain []*model.Blockchain

	result := repository.IDatabaseHandler.GetClient().Preload("WalletAddresses").Model(&model.Blockchain{}).Scopes()
	result = result.Find(&blockchain)

	if result.Error != nil {
		return nil, fmt.Errorf("blockchain get list failed: %s", result.Error.Error())
	}

	return blockchain, nil

}

func (repository *BlockchainRepository) GetById(id uint) (*model.Blockchain, error) {
	var blockchain model.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&blockchain, "id = ?", id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("blockchain get by id failed: %s", result.Error.Error())
	}

	return &blockchain, nil
}

func (repository *BlockchainRepository) GetByUuid(uuid *uuid.UUID) (*model.Blockchain, error) {
	var blockchain model.Blockchain

	result := repository.IDatabaseHandler.GetClient().First(&blockchain, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("blockchain get by uuid failed: %s", result.Error.Error())
	}

	return &blockchain, nil
}

func (repository *BlockchainRepository) GetByName(name string) (*model.Blockchain, error) {
	var blockchain model.Blockchain

	result := repository.IDatabaseHandler.GetClient().Where("is_active = ? ", true).First(&blockchain, "name = ?", name)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("blockchain get by name failed: %s", result.Error.Error())
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
	result = result.Count(&count)

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
			return nil, fmt.Errorf("blockchain with UUID %s not found", uuid)
		}

		return nil, fmt.Errorf("failed to retrieve blockchain with UUID %s: %s", uuid, result.Error)
	}

	if err := repository.IDatabaseHandler.GetClient().Model(&existingBlockchain).Updates(req).Error; err != nil {
		return nil, fmt.Errorf("failed to update blockchain: %s", err)
	}

	return &existingBlockchain, nil
}
