package repositories

import (
	"athena/src/database"
	blockchain_model "athena/src/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IWalletAddressRepository interface {
	GetList() ([]*blockchain_model.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*blockchain_model.WalletAddress, error)
	Create(blockchain *blockchain_model.WalletAddress) (*blockchain_model.WalletAddress, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *blockchain_model.WalletAddress) (*blockchain_model.WalletAddress, error)
	GetActiveList() ([]*blockchain_model.WalletAddress, error)
	GetWalletAddress(blockchainName string, number int) ([]string, error)
}

type WalletAddressRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *WalletAddressRepository) GetWalletAddress(blockchainName string, count int) ([]string, error) {
	var blockchain blockchain_model.Blockchain
	var walletAddressesName []string

	// First, find the blockchain by name and preload related wallet addresses
	if err := repository.IDatabaseHandler.GetClient().
		Preload("WalletAddresses", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ? AND allocated_at IS NULL", true).Limit(count)
		}).
		Where("name = ?", blockchainName).
		First(&blockchain).Error; err != nil {
		return nil, err
	}

	// If no wallet addresses are found, return an error
	if len(blockchain.WalletAddresses) == 0 {
		return nil, errors.New("no wallet addresses found")
	}

	// Update the found wallet addresses to allocated
	for _, walletAddress := range blockchain.WalletAddresses {
		walletAddress.AllocatedAt = time.Now()
		if err := repository.IDatabaseHandler.GetClient().Save(&walletAddress).Error; err != nil {
			return nil, err
		}
	}

	// Collect wallet addresses for return
	for _, walletAddress := range blockchain.WalletAddresses {
		walletAddressesName = append(walletAddressesName, walletAddress.WalletAddress)
	}

	return walletAddressesName, nil
}

func (repository *WalletAddressRepository) GetActiveList() ([]*blockchain_model.WalletAddress, error) {
	var walletAddress []*blockchain_model.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Where("is_active = ?", true).Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) GetList() ([]*blockchain_model.WalletAddress, error) {
	var walletAddress []*blockchain_model.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) GetByUuid(uuid *uuid.UUID) (*blockchain_model.WalletAddress, error) {
	var walletAddress blockchain_model.WalletAddress

	result := repository.IDatabaseHandler.GetClient().First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("walletAddress get by uuid failed: %s", result.Error.Error())
	}

	return &walletAddress, nil
}

func (repository *WalletAddressRepository) Create(walletAddress *blockchain_model.WalletAddress) (*blockchain_model.WalletAddress, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&walletAddress)
	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress creation failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) GetCount() (int64, error) {
	var count int64

	result := repository.IDatabaseHandler.GetClient().Model(&blockchain_model.WalletAddress{}).Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("walletAddresses get count failed: %s", result.Error.Error())
	}

	return count, nil

}

func (repository *WalletAddressRepository) Delete(uuid *uuid.UUID) error {
	var walletAddress blockchain_model.WalletAddress

	result := repository.IDatabaseHandler.GetClient().First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("walletAddress get by uuid failed: %s", result.Error.Error())
	}

	if err := repository.IDatabaseHandler.GetClient().Delete(&walletAddress).Error; err != nil {
		return fmt.Errorf("failed to delete walletAddress: %s", err)
	}

	return nil
}

func (repository *WalletAddressRepository) Update(uuid *uuid.UUID, req *blockchain_model.WalletAddress) (*blockchain_model.WalletAddress, error) {
	var walletAddress blockchain_model.WalletAddress

	result := repository.IDatabaseHandler.GetClient().First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("walletAddress with UUID %s not found", uuid)
		}

		return nil, fmt.Errorf("failed to retrieve walletAddress with UUID %s: %s", uuid, result.Error)
	}

	if err := repository.IDatabaseHandler.GetClient().Model(&walletAddress).Updates(req).Error; err != nil {
		return nil, fmt.Errorf("failed to update walletAddress: %s", err)
	}

	return &walletAddress, nil
}
