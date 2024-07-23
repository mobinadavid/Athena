package repository

import (
	"athena/src/database"
	"athena/src/services/wallet-address/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IWalletAddressRepository interface {
	GetList() ([]*model.WalletAddress, error)
	GetByUuid(uuid *uuid.UUID) (*model.WalletAddress, error)
	Create(blockchain *model.WalletAddress) (*model.WalletAddress, error)
	GetCount() (int64, error)
	Delete(uuid *uuid.UUID) error
	Update(uuid *uuid.UUID, req *model.WalletAddress) (*model.WalletAddress, error)
	GetActiveList() ([]*model.WalletAddress, error)
}

type WalletAddressRepository struct {
	IDatabaseHandler *database.Database
}

func (repository *WalletAddressRepository) GetActiveList() ([]*model.WalletAddress, error) {
	var walletAddress []*model.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Where("is_active = ?", true).Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}
func (repository *WalletAddressRepository) GetList() ([]*model.WalletAddress, error) {
	var walletAddress []*model.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) GetByUuid(uuid *uuid.UUID) (*model.WalletAddress, error) {
	var walletAddress model.WalletAddress

	result := repository.IDatabaseHandler.GetClient().First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("walletAddress get by uuid failed: %s", result.Error.Error())
	}

	return &walletAddress, nil
}

func (repository *WalletAddressRepository) Create(walletAddress *model.WalletAddress) (*model.WalletAddress, error) {
	result := repository.IDatabaseHandler.GetClient().Create(&walletAddress)
	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress creation failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) GetCount() (int64, error) {
	var count int64

	result := repository.IDatabaseHandler.GetClient().Model(&model.WalletAddress{})

	if result.Error != nil {
		return 0, fmt.Errorf("walletAddresses get count failed: %s", result.Error.Error())
	}

	return count, nil
}

func (repository *WalletAddressRepository) Delete(uuid *uuid.UUID) error {
	var walletAddress model.WalletAddress

	result := repository.IDatabaseHandler.GetClient().First(&walletAddress, "uuid = ?", uuid)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("walletAddress get by uuid failed: %s", result.Error.Error())
	}

	if err := repository.IDatabaseHandler.GetClient().Delete(&walletAddress).Error; err != nil {
		return fmt.Errorf("failed to delete walletAddress: %s", err)
	}

	return nil
}

func (repository *WalletAddressRepository) Update(uuid *uuid.UUID, req *model.WalletAddress) (*model.WalletAddress, error) {
	var walletAddress model.WalletAddress

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
