package repositories

import (
	"athena/src/api/errs"
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	IsAllocated(address *models.WalletAddress) (bool, error)
	GetUnallocatedWalletAddress(count int, blockchainId uint) ([]*models.WalletAddress, error)
	UpdateWalletAddressToAllocated(walletAddresses []*models.WalletAddress, userID *uint, paymentRequestID *uint) error
	ReleaseWalletAddresses(walletAddresses []*models.WalletAddress) error
	GetAllocatedByUser(userID uint) ([]*models.WalletAddress, error)
	AllocateAtomic(count int, blockchainId uint, userID *uint, paymentRequestID *uint) ([]*models.WalletAddress, error)
	PoolStats() (total, allocated, free int64, err error)
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

func (repository *WalletAddressRepository) UpdateWalletAddressToAllocated(walletAddresses []*models.WalletAddress, userID *uint, paymentRequestID *uint) error {
	now := time.Now()
	for _, walletAddress := range walletAddresses {
		walletAddress.AllocatedAt = &now
		walletAddress.AllocatedToUserID = userID
		walletAddress.PaymentRequestID = paymentRequestID
		if err := repository.IDatabaseHandler.GetClient().Save(walletAddress).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repository *WalletAddressRepository) ReleaseWalletAddresses(walletAddresses []*models.WalletAddress) error {
	for _, walletAddress := range walletAddresses {
		if err := repository.IDatabaseHandler.GetClient().Model(walletAddress).Updates(map[string]interface{}{
			"allocated_at":         nil,
			"allocated_to_user_id": nil,
			"payment_request_id":   nil,
		}).Error; err != nil {
			return err
		}
		walletAddress.AllocatedAt = nil
		walletAddress.AllocatedToUserID = nil
		walletAddress.PaymentRequestID = nil
	}
	return nil
}

func (repository *WalletAddressRepository) GetAllocatedByUser(userID uint) ([]*models.WalletAddress, error) {
	var walletAddress []*models.WalletAddress
	result := repository.IDatabaseHandler.GetClient().
		Scopes(scopes.IsAllocated()).
		Preload("Blockchain").
		Where("allocated_to_user_id = ?", userID).
		Find(&walletAddress)
	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}
	return walletAddress, nil
}

func (repository *WalletAddressRepository) IsAllocated(walletAddress *models.WalletAddress) (bool, error) {
	var count int64
	err := repository.IDatabaseHandler.GetClient().
		Model(&models.WalletAddress{}).
		Where("wallet_address = ?", walletAddress.WalletAddress).Scopes(scopes.IsAllocated()).Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *WalletAddressRepository) GetUnallocatedWalletAddress(count int, blockchainId uint) ([]*models.WalletAddress, error) {
	var walletAddresses []*models.WalletAddress
	// Use IsAllocated to filter records
	result := repository.IDatabaseHandler.GetClient().
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Scopes(scopes.IsNotAllocated()).
		Preload("Blockchain").
		Where("blockchain_id = ?", blockchainId).
		Limit(count)
	if err := result.Find(&walletAddresses).Error; err != nil {
		return nil, err
	}

	return walletAddresses, nil
}

func (repository *WalletAddressRepository) GetAllocatedList() ([]*models.WalletAddress, error) {
	var walletAddress []*models.WalletAddress

	result := repository.IDatabaseHandler.GetClient()
	result = result.Scopes(scopes.IsAllocated()).
		Preload("Blockchain").
		Preload("PaymentRequest").
		Find(&walletAddress)

	if result.Error != nil {
		return nil, fmt.Errorf("walletAddress get list failed: %s", result.Error.Error())
	}

	return walletAddress, nil
}

func (repository *WalletAddressRepository) PoolStats() (int64, int64, int64, error) {
	var total, allocated int64
	if err := repository.IDatabaseHandler.GetClient().Model(&models.WalletAddress{}).Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}
	if err := repository.IDatabaseHandler.GetClient().Model(&models.WalletAddress{}).Where("allocated_at is not null").Count(&allocated).Error; err != nil {
		return 0, 0, 0, err
	}
	return total, allocated, total - allocated, nil
}

func (repository *WalletAddressRepository) AllocateAtomic(count int, blockchainId uint, userID *uint, paymentRequestID *uint) ([]*models.WalletAddress, error) {
	var walletAddresses []*models.WalletAddress
	err := repository.IDatabaseHandler.GetClient().Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Scopes(scopes.IsNotAllocated()).
			Preload("Blockchain").
			Where("blockchain_id = ?", blockchainId).
			Limit(count).
			Find(&walletAddresses).Error; err != nil {
			return err
		}
		if len(walletAddresses) != count {
			return errs.ErrNotEnoughWalletAddresses
		}
		now := time.Now()
		for _, walletAddress := range walletAddresses {
			walletAddress.AllocatedAt = &now
			walletAddress.AllocatedToUserID = userID
			walletAddress.PaymentRequestID = paymentRequestID
			if err := tx.Save(walletAddress).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return walletAddresses, nil
}
