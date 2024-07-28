package repository

import (
	"athena/src/database"
	"athena/src/services/blockchain-explorer/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBlockchainExplorerRepository interface {
	GetList() ([]*model.BlockchainExplorer, error)
	GetByUuid(uuid *uuid.UUID) (*model.BlockchainExplorer, error)
	Create(explorer *model.BlockchainExplorer) (*model.BlockchainExplorer, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *model.BlockchainExplorer) (*model.BlockchainExplorer, error)
	GetExplorerByBlockchainID(blockchainID uint) (*model.BlockchainExplorer, error)
}

type BlockchainExplorerRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *BlockchainExplorerRepository) Create(blockchainExplorer *model.BlockchainExplorer) (*model.BlockchainExplorer, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&blockchainExplorer)
	if result.Error != nil {
		return nil, fmt.Errorf("blockchainExplorer creation failed: %s", result.Error.Error())
	}

	return blockchainExplorer, nil
}

func (repository *BlockchainExplorerRepository) GetList() ([]*model.BlockchainExplorer, error) {
	var blockchainExplorers []*model.BlockchainExplorer
	result := repository.IDatabaseHandler.GetClient().Preload("Blockchains").Model(&model.BlockchainExplorer{})

	result = result.Scopes()

	result = result.Find(&blockchainExplorers)

	if result.Error != nil {
		return nil, fmt.Errorf("blockchainExplorers get list failed: %s", result.Error.Error())
	}

	return blockchainExplorers, nil
}

func (repository *BlockchainExplorerRepository) GetByUuid(uuid *uuid.UUID) (*model.BlockchainExplorer, error) {
	var blockchainExplorer model.BlockchainExplorer

	// Preload blockchains
	result := repository.IDatabaseHandler.GetClient().Preload("Blockchains").First(&blockchainExplorer, "uuid = ?", uuid)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blockchainExplorer not found for uuid: %s", uuid)
		}
		return nil, fmt.Errorf("failed to get blockchainExplorer by uuid: %s", result.Error.Error())
	}

	return &blockchainExplorer, nil
}

func (repository *BlockchainExplorerRepository) GetCount() (int64, error) {
	var count int64
	result := repository.IDatabaseHandler.GetClient().Model(&model.BlockchainExplorer{})

	result = result.Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("blockchainExplorers get count failed: %s", result.Error.Error())
	}

	return count, nil
}

func (repository *BlockchainExplorerRepository) Delete(uuid *uuid.UUID) error {
	var blockchainExplorer model.BlockchainExplorer

	result := repository.IDatabaseHandler.GetClient().Preload("Blockchains").First(&blockchainExplorer, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("blockchainExplorer get by uuid failed: %s", result.Error.Error())
	}

	// Delete associated blockchains
	if len(blockchainExplorer.Blockchains) > 0 {
		if err := repository.IDatabaseHandler.GetClient().Model(&blockchainExplorer).Association("Blockchains").Clear(); err != nil {
			return fmt.Errorf("failed to clear blockchain association for blockchainExplorer: %s", err)
		}
	}

	if err := repository.IDatabaseHandler.GetClient().Delete(&blockchainExplorer).Error; err != nil {
		return fmt.Errorf("failed to delete blockchainExplorer: %s", err)
	}

	return nil
}

func (repository *BlockchainExplorerRepository) Update(uuid *uuid.UUID, blockchainExplorer *model.BlockchainExplorer) (*model.BlockchainExplorer, error) {
	var existing model.BlockchainExplorer
	var result = repository.IDatabaseHandler.GetClient().Preload("Blockchains").First(&existing, "uuid = ?", uuid)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("blockchainExplorer with UUID %s not found", uuid)
		}
		return nil, fmt.Errorf("failed to retrieve blockchainExplorer with UUID %s: %s", uuid, result.Error)
	}

	if err := repository.IDatabaseHandler.GetClient().Model(&existing).Association("Blockchains").Clear(); err != nil {
		return nil, fmt.Errorf("failed to clear blockchains: %s", err)
	}

	if err := repository.IDatabaseHandler.GetClient().Model(&existing).Association("Blockchains").Append(blockchainExplorer.Blockchains); err != nil {
		return nil, fmt.Errorf("failed to update blockchains: %s", err)
	}

	if err := repository.IDatabaseHandler.GetClient().Session(&gorm.Session{FullSaveAssociations: true}).Model(&existing).Updates(blockchainExplorer).Error; err != nil {
		return nil, fmt.Errorf("failed to update blockchainExplorer: %s", err)
	}

	return &existing, nil
}

func (repository *BlockchainExplorerRepository) GetExplorerByBlockchainID(blockchainID uint) (*model.BlockchainExplorer, error) {
	var explorer model.BlockchainExplorer

	result := repository.IDatabaseHandler.GetClient().Joins("JOIN blockchain_explorer_mappings ON blockchain_explorer_mappings.blockchain_explorer_id = blockchain_explorers.id").Where("blockchain_explorer_mappings.blockchain_id = ? AND blockchain_explorers.is_active = ? AND blockchain_explorers.is_default = ?", blockchainID, true, true).First(&explorer)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active default explorer for blockchain ID %d: %w", blockchainID, result.Error)
	}

	return &explorer, nil

}
