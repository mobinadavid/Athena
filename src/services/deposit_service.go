package services

import (
	"athena/src/api/errs"
	"athena/src/config"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/logger"
	"athena/src/repositories"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type IDepositService interface {
	TransactionsExist(txHash string) (bool, error)
	Create(deposit *models.Deposits) error
	HandleDeposits() error
	GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error)
	GetByUuid(userID *uint, id *uuid.UUID) (*models.Deposits, error)
	Dashboard(userID *uint) (map[string]interface{}, error)
}

type DepositService struct {
	IDepositRepository       repositories.IDepositRepository
	IWalletAddressService    IWalletAddressService
	PaymentRequestRepository repositories.IPaymentRequestRepository
	NotificationService      INotificationService
}

func (service *DepositService) TransactionsExist(txHash string) (bool, error) {
	return service.IDepositRepository.TransactionsExist(txHash)
}

func (service *DepositService) Create(deposit *models.Deposits) error {
	return service.IDepositRepository.Create(deposit)
}

func (service *DepositService) GetList(params *scopes.QueryBuilderModel) (*scopes.PaginatedModel, error) {
	items, count, err := service.IDepositRepository.GetList(params)
	if err != nil {
		return nil, err
	}
	return paginatedResult(params, items, count), nil
}

func (service *DepositService) GetByUuid(userID *uint, id *uuid.UUID) (*models.Deposits, error) {
	deposit, err := service.IDepositRepository.GetByUuid(id)
	if err != nil {
		return nil, errs.RecordNotFound
	}
	if userID != nil && (deposit.UserID == nil || *deposit.UserID != *userID) {
		return nil, errs.ErrForbiddenResource
	}
	return deposit, nil
}

func (service *DepositService) Dashboard(userID *uint) (map[string]interface{}, error) {
	statusCounts, err := service.IDepositRepository.CountByStatus(userID)
	if err != nil {
		return nil, err
	}
	totalAmount, err := service.IDepositRepository.SumAmount(userID)
	if err != nil {
		return nil, err
	}
	series, err := service.IDepositRepository.DashboardSeries(userID, 14)
	if err != nil {
		return nil, err
	}
	byBlockchain, err := service.IDepositRepository.DashboardByBlockchain(userID)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"pending_payments":       statusCounts[models.PaymentStatusPending],
		"confirmed_payments":     statusCounts[models.PaymentStatusConfirmed],
		"expired_payments":       statusCounts[models.PaymentStatusExpired],
		"total_amount":           totalAmount,
		"deposits_over_time":     series,
		"deposits_by_blockchain": byBlockchain,
	}

	if userID != nil {
		wallets, err := service.IWalletAddressService.GetAllocatedByUser(*userID)
		if err != nil {
			return nil, err
		}
		unread := int64(0)
		if service.NotificationService != nil {
			unread, _ = service.NotificationService.UnreadCount(*userID)
		}
		data["allocated_wallets"] = len(wallets)
		data["unread_notifications"] = unread
	} else {
		allocated, err := service.IWalletAddressService.GetAllocatedList()
		if err == nil {
			if wallets, ok := allocated.Items.(*[]*models.WalletAddress); ok {
				data["allocated_wallets"] = len(*wallets)
			}
		}
	}

	return data, nil
}

func (service *DepositService) HandleDeposits() error {
	if err := service.expirePendingAllocations(); err != nil {
		logger.GetInstance().Sugar().Warnf("failed to expire pending allocations: %v", err)
	}

	allocatedList, err := service.IWalletAddressService.GetAllocatedList()
	if err != nil {
		return err
	}

	walletAddresses := allocatedList.Items.(*[]*models.WalletAddress)
	for _, walletAddress := range *walletAddresses {
		if err := service.trackWallet(walletAddress); err != nil {
			logger.GetInstance().Sugar().Warnf("failed to track wallet %s: %v", walletAddress.WalletAddress, err)
		}
	}
	return nil
}

func (service *DepositService) expirePendingAllocations() error {
	expired, err := service.PaymentRequestRepository.GetExpiredPending()
	if err != nil {
		return err
	}
	for _, request := range expired {
		if len(request.WalletAddresses) > 0 {
			if err := service.IWalletAddressService.ReleaseWalletAddresses(request.WalletAddresses); err != nil {
				logger.GetInstance().Sugar().Warnf("failed to release expired wallets for payment %s: %v", request.UUID, err)
				continue
			}
		}
		request.Status = models.PaymentStatusExpired
		if _, err := service.PaymentRequestRepository.Update(request); err != nil {
			continue
		}
		_ = service.notifyUser(request.UserID, "payment_expired", "پرداخت منقضی شد", "مهلت پرداخت به پایان رسید و آدرس‌ها آزاد شدند.", map[string]interface{}{
			"payment_request_uuid": request.UUID.String(),
		})
	}
	return nil
}

func (service *DepositService) trackWallet(walletAddress *models.WalletAddress) error {
	transactions, err := service.IWalletAddressService.GetTransactionsList(walletAddress)
	if err != nil {
		return fmt.Errorf("failed to get transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
	}

	filteredTxs, err := service.IWalletAddressService.FilterTransactions(transactions, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to filter transactions for wallet address %s: %w", walletAddress.WalletAddress, err)
	}

	paid := false
	for _, tx := range filteredTxs {
		exists, err := service.IDepositRepository.TransactionsExist(tx.Hash)
		if err != nil {
			return fmt.Errorf("failed to check transaction existence: %w", err)
		}
		if exists {
			paid = true
			continue
		}

		deposit := depositFromTransaction(tx, walletAddress)
		if err := service.IDepositRepository.Create(deposit); err != nil {
			logger.GetInstance().Sugar().Warnf("failed to add transaction to deposits: %v", err)
			continue
		}
		paid = true

		if walletAddress.WebhookURL != "" {
			if err := sendToWebhook(tx, walletAddress.WebhookURL); err != nil {
				logger.GetInstance().Sugar().Warnf("failed to send transaction to wallet webhook: %v", err)
			}
		}
		if err := notifyFinancialSystem(deposit); err != nil {
			logger.GetInstance().Sugar().Warnf("failed to notify financial system: %v", err)
		}
		if walletAddress.AllocatedToUserID != nil {
			_ = service.notifyUser(*walletAddress.AllocatedToUserID, "payment_confirmed", "پرداخت تایید شد", "تراکنش شما در شبکه بلاکچین تایید شد.", map[string]interface{}{
				"transaction_hash":     tx.Hash,
				"amount":               tx.Value,
				"blockchain":           tx.BlockChain,
				"wallet_address":       walletAddress.WalletAddress,
				"payment_request_uuid": paymentRequestUUID(walletAddress),
			})
		}
	}

	if paid {
		if walletAddress.PaymentRequestID != nil {
			service.tryConfirmPayment(*walletAddress.PaymentRequestID)
		}
		if err := service.IWalletAddressService.ReleaseWalletAddresses([]*models.WalletAddress{walletAddress}); err != nil {
			return err
		}
	}
	return nil
}

func (service *DepositService) tryConfirmPayment(paymentRequestID uint) {
	request, err := service.PaymentRequestRepository.GetByID(paymentRequestID)
	if err != nil || request.Status == models.PaymentStatusConfirmed {
		return
	}

	if request.ExpectedAmount != nil {
		deposits, err := service.IDepositRepository.GetByPaymentRequest(request.ID)
		if err != nil {
			return
		}
		var total float64
		for _, deposit := range deposits {
			total += deposit.Amount
		}
		if total < *request.ExpectedAmount {
			return
		}
	}

	service.markPaymentConfirmed(paymentRequestID)
}

func (service *DepositService) markPaymentConfirmed(paymentRequestID uint) {
	request, err := service.PaymentRequestRepository.GetByID(paymentRequestID)
	if err != nil {
		return
	}
	if request.Status == models.PaymentStatusConfirmed {
		return
	}
	now := time.Now()
	request.Status = models.PaymentStatusConfirmed
	request.ConfirmedAt = &now
	_, _ = service.PaymentRequestRepository.Update(request)
}

func (service *DepositService) notifyUser(userID uint, notificationType, title, body string, data map[string]interface{}) error {
	if service.NotificationService == nil {
		return nil
	}
	_, err := service.NotificationService.Create(userID, notificationType, title, body, data)
	return err
}

func depositFromTransaction(tx *models.TransactionResponse, walletAddress *models.WalletAddress) *models.Deposits {
	toAddress := walletAddress.WalletAddress
	if len(tx.ToAddresses) > 0 {
		toAddress = tx.ToAddresses[0]
	}
	timestamp := tx.Timestamp
	now := time.Now()
	return &models.Deposits{
		TransactionHash:  tx.Hash,
		FromAddress:      tx.From,
		ToAddress:        toAddress,
		Amount:           tx.Value,
		Fee:              tx.Fee,
		Confirmations:    tx.Confirmations,
		Status:           models.DepositStatusConfirmed,
		BlockNumber:      tx.BlockNumber,
		PaidAt:           &timestamp,
		BlockchainID:     &walletAddress.BlockchainID,
		WalletAddressID:  &walletAddress.ID,
		PaymentRequestID: walletAddress.PaymentRequestID,
		UserID:           walletAddress.AllocatedToUserID,
		NotifiedAt:       &now,
	}
}

func paymentRequestUUID(walletAddress *models.WalletAddress) string {
	if walletAddress.PaymentRequest != nil {
		return walletAddress.PaymentRequest.UUID.String()
	}
	return ""
}

func notifyFinancialSystem(deposit *models.Deposits) error {
	webhookURL := config.GetInstance().Get("FINANCIAL_SYSTEM_WEBHOOK_URL")
	if webhookURL == "" {
		return nil
	}
	payload, err := json.Marshal(deposit)
	if err != nil {
		return err
	}
	client := resty.New().SetTimeout(10 * time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := client.R().SetHeader("Content-Type", "application/json").SetBody(payload).Post(webhookURL)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 300 {
		return fmt.Errorf("financial system webhook returned status %d", resp.StatusCode())
	}
	return nil
}

func sendToWebhook(tx *models.TransactionResponse, webhookUrl string) error {
	txData, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	configs := config.GetInstance()
	maxRetry, _ := strconv.Atoi(configs.Get("SEND_TO_WEBHOOK_MAX_RETRY"))

	client := resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	webhookUrl = webhookUrl + "/" + tx.Hash
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(txData).
		Post(webhookUrl)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("webhook returned non-200 status: %d", resp.StatusCode())
	}
	return nil
}
