package repositories

import (
	"athena/src/database"
	"athena/src/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IWalletAddressRepository interface {
	GetList() ([]*models.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*models.WalletAddress, error)
	Create(blockchain *models.WalletAddress) (*models.WalletAddress, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *models.WalletAddress) (*models.WalletAddress, error)
	GetActiveList() ([]*models.WalletAddress, error)
	GetWalletAddress(blockchainName string, number int) ([]string, error)
}

type WalletAddressRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *WalletAddressRepository) GetWalletAddress(blockchainName string, count int) ([]string, error) {
	var walletAddresses []*models.WalletAddress
	var walletAddressesName []string

	if err := repository.IDatabaseHandler.GetClient().
		Preload("Blockchain").
		Where("wallet_addresses.is_active = ? AND wallet_addresses.allocated_at IS NULL", true).
		Limit(count).
		Find(&walletAddresses).Error; err != nil {
		return nil, err
	}

	// Filter wallet addresses based on the blockchain name
	var filteredWalletAddresses []*models.WalletAddress
	for _, walletAddress := range walletAddresses {
		if walletAddress.Blockchain.Name == blockchainName {
			filteredWalletAddresses = append(filteredWalletAddresses, walletAddress)
		}
	}

	// Check if any addresses are found
	if len(filteredWalletAddresses) == 0 {
		return nil, errors.New("no wallet addresses found")
	}

	// Update the found wallet addresses to allocated
	for _, walletAddress := range filteredWalletAddresses {
		walletAddress.AllocatedAt = time.Now()
		if err := repository.IDatabaseHandler.GetClient().Save(walletAddress).Error; err != nil {
			return nil, err
		}
	}

	// Collect wallet address names
	for _, walletAddress := range filteredWalletAddresses {
		walletAddressesName = append(walletAddressesName, walletAddress.WalletAddress)
	}

	return walletAddressesName, nil
}

func (repository *WalletAddressRepository) GetActiveList() ([]*models.WalletAddress, error) {
	var walletAddress []*models.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Where("is_active = ?", true).Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) GetList() ([]*models.WalletAddress, error) {
	var walletAddress []*models.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Preload("Blockchain").Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
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

func (repository *WalletAddressRepository) GetCount() (int64, error) {
	var count int64

	result := repository.IDatabaseHandler.GetClient().Model(&models.WalletAddress{}).Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("walletAddresses get count failed: %s", result.Error.Error())
	}

	return count, nil

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
