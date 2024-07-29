package repositories

import (
	"athena/src/database"
	"athena/src/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBlockchainExplorerRepository interface {
	GetList() ([]*models.BlockchainExplorer, error)
	GetByUuid(uuid *uuid.UUID) (*models.BlockchainExplorer, error)
	Create(explorer *models.BlockchainExplorer) (*models.BlockchainExplorer, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *models.BlockchainExplorer) (*models.BlockchainExplorer, error)
	GetExplorerByBlockchain(blockchain *models.Blockchain) (*models.BlockchainExplorer, error)
}

type BlockchainExplorerRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *BlockchainExplorerRepository) Create(blockchainExplorer *models.BlockchainExplorer) (*models.BlockchainExplorer, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&blockchainExplorer)
	if result.Error != nil {
		return nil, fmt.Errorf("blockchainExplorer creation failed: %s", result.Error.Error())
	}

	return blockchainExplorer, nil
}

func (repository *BlockchainExplorerRepository) GetList() ([]*models.BlockchainExplorer, error) {
	var blockchainExplorers []*models.BlockchainExplorer
	result := repository.IDatabaseHandler.GetClient().Preload("Blockchains").Model(&models.BlockchainExplorer{})

	result = result.Scopes()

	result = result.Find(&blockchainExplorers)

	if result.Error != nil {
		return nil, fmt.Errorf("blockchainExplorers get list failed: %s", result.Error.Error())
	}

	return blockchainExplorers, nil
}

func (repository *BlockchainExplorerRepository) GetByUuid(uuid *uuid.UUID) (*models.BlockchainExplorer, error) {
	var blockchainExplorer models.BlockchainExplorer

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
	result := repository.IDatabaseHandler.GetClient().Model(&models.BlockchainExplorer{})

	result = result.Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("blockchainExplorers get count failed: %s", result.Error.Error())
	}

	return count, nil
}

func (repository *BlockchainExplorerRepository) Delete(uuid *uuid.UUID) error {
	var blockchainExplorer models.BlockchainExplorer

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

func (repository *BlockchainExplorerRepository) Update(uuid *uuid.UUID, blockchainExplorer *models.BlockchainExplorer) (*models.BlockchainExplorer, error) {
	var existing models.BlockchainExplorer
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

func (repository *BlockchainExplorerRepository) GetExplorerByBlockchain(blockchain *models.Blockchain) (*models.BlockchainExplorer, error) {
	var explorers []models.BlockchainExplorer

	// Preload blockchains and apply conditions
	result := repository.IDatabaseHandler.GetClient().
		Preload("Blockchains").
		Where("blockchain_explorers.is_active = ? AND blockchain_explorers.is_default = ?", true, true).
		Find(&explorers)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active default explorer: %w", result.Error)
	}

	// Filter explorers based on the blockchain ID
	for _, explorer := range explorers {
		for _, b := range explorer.Blockchains {
			if b.ID == blockchain.ID {
				return &explorer, nil
			}
		}
	}

	return nil, fmt.Errorf("no active default explorer found for blockchain ID %d", blockchain.ID)

}
