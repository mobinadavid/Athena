package services

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests"
	"athena/src/config"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/repositories"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type IPaymentService interface {
	Create(userID uint, request *requests.AllocateWalletAddress) (*models.PaymentRequest, error)
	GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error)
	GetByUuid(userID *uint, id *uuid.UUID) (*models.PaymentRequest, error)
	GetTransactions(userID *uint, id *uuid.UUID, page, limit uint) (*scopes.PaginatedModel, error)
}

type PaymentService struct {
	PaymentRequestRepository repositories.IPaymentRequestRepository
	WalletAddressService     IWalletAddressService
	WalletAddressRepository  repositories.IWalletAddressRepository
	BlockchainService        IBlockChainService
	DepositRepository        repositories.IDepositRepository
}

func (service *PaymentService) Create(userID uint, request *requests.AllocateWalletAddress) (*models.PaymentRequest, error) {
	if request.Count <= 0 {
		return nil, errs.ErrInvalidWalletCount
	}

	blockchain, err := service.BlockchainService.GetByName(request.Blockchain)
	if err != nil {
		return nil, err
	}

	ttlMinutes := envIntFallback("ALLOCATION_TTL_MINUTES", 30)
	expiresAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)
	paymentRequest := &models.PaymentRequest{
		UserID:         userID,
		BlockchainID:   blockchain.ID,
		RequestedCount: request.Count,
		ExpectedAmount: request.ExpectedAmount,
		Status:         models.PaymentStatusPending,
		ExpiresAt:      &expiresAt,
	}

	created, err := service.PaymentRequestRepository.Create(paymentRequest)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	wallets, err := service.WalletAddressRepository.AllocateAtomic(request.Count, blockchain.ID, &userID, &created.ID)
	if err != nil {
		_ = service.PaymentRequestRepository.Delete(created.ID)
		return nil, err
	}

	created.Blockchain = blockchain
	created.WalletAddresses = wallets
	return created, nil
}

func (service *PaymentService) GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error) {
	items, count, err := service.PaymentRequestRepository.GetList(params)
	if err != nil {
		return nil, err
	}
	return paginatedResult(params, items, count), nil
}

func (service *PaymentService) GetByUuid(userID *uint, id *uuid.UUID) (*models.PaymentRequest, error) {
	request, err := service.PaymentRequestRepository.GetByUuid(id, "Blockchain", "WalletAddresses", "WalletAddresses.Blockchain", "Deposits")
	if err != nil {
		return nil, errs.RecordNotFound
	}
	if userID != nil && request.UserID != *userID {
		return nil, errs.ErrForbiddenResource
	}
	return request, nil
}

func (service *PaymentService) GetTransactions(userID *uint, id *uuid.UUID, page, limit uint) (*scopes.PaginatedModel, error) {
	request, err := service.GetByUuid(userID, id)
	if err != nil {
		return nil, err
	}

	deposits, err := service.DepositRepository.GetByPaymentRequest(request.ID)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	start := (page - 1) * limit
	end := start + limit
	if start > uint(len(deposits)) {
		start = uint(len(deposits))
	}
	if end > uint(len(deposits)) {
		end = uint(len(deposits))
	}

	paged := deposits[start:end]
	total := int64(len(deposits))
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	return &scopes.PaginatedModel{
		Limit:       limit,
		CurrentPage: page,
		TotalItems:  total,
		TotalPages:  totalPages,
		Items:       &paged,
	}, nil
}

func envIntFallback(key string, fallback int) int {
	value, err := strconv.Atoi(config.GetInstance().Get(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
