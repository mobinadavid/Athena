package providers

import "athena/src/services"

func ProvideOTPService() *services.OTPService {
	return &services.OTPService{
		//	NotificationService: notificationService,
	}
}
