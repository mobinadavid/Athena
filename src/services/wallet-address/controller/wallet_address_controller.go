package controller

import "athena/src/services/wallet-address/service"

type WalletAddressController struct {
	IWalletAddressService *service.WalletAddressService
}
