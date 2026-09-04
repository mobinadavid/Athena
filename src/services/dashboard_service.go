package services

import (
	"athena/src/api/errs"
	"athena/src/repositories"
)

type IDashboardService interface {
	UserSummary(userID uint) (map[string]interface{}, error)
	AdminSummary() (map[string]interface{}, error)
	DepositsOverTime(userID *uint, days int) ([]map[string]interface{}, error)
	DepositsByBlockchain(userID *uint) ([]map[string]interface{}, error)
}

type DashboardService struct {
	DepositRepository       repositories.IDepositRepository
	WalletAddressRepository repositories.IWalletAddressRepository
	NotificationRepository  repositories.INotificationRepository
}

func (service *DashboardService) UserSummary(userID uint) (map[string]interface{}, error) {
	userIDPtr := &userID
	series, err := service.DepositRepository.DashboardSeries(userIDPtr, 14)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	byBlockchain, err := service.DepositRepository.DashboardByBlockchain(userIDPtr)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	statusCounts, err := service.DepositRepository.CountByStatus(userIDPtr)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	amount, err := service.DepositRepository.SumAmount(userIDPtr)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	unread, err := service.NotificationRepository.UnreadCount(userID)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	wallets, err := service.WalletAddressRepository.GetAllocatedByUser(userID)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return map[string]interface{}{
		"payment_requests_by_status": statusCounts,
		"allocated_wallets":          len(wallets),
		"total_received_amount":      amount,
		"unread_notifications":       unread,
		"deposits_over_time":         series,
		"deposits_by_blockchain":     byBlockchain,
	}, nil
}

func (service *DashboardService) AdminSummary() (map[string]interface{}, error) {
	series, err := service.DepositRepository.DashboardSeries(nil, 30)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	byBlockchain, err := service.DepositRepository.DashboardByBlockchain(nil)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	statusCounts, err := service.DepositRepository.CountByStatus(nil)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	amount, err := service.DepositRepository.SumAmount(nil)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	total, allocated, free, err := service.WalletAddressRepository.PoolStats()
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return map[string]interface{}{
		"payment_requests_by_status": statusCounts,
		"total_received_amount":      amount,
		"wallet_pool": map[string]int64{
			"total":     total,
			"allocated": allocated,
			"free":      free,
		},
		"deposits_over_time":     series,
		"deposits_by_blockchain": byBlockchain,
	}, nil
}

func (service *DashboardService) DepositsOverTime(userID *uint, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 14
	}
	if days > 90 {
		days = 90
	}
	series, err := service.DepositRepository.DashboardSeries(userID, days)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	return series, nil
}

func (service *DashboardService) DepositsByBlockchain(userID *uint) ([]map[string]interface{}, error) {
	items, err := service.DepositRepository.DashboardByBlockchain(userID)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}
	return items, nil
}
