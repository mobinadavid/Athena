package service

import "athena/src/services/wallet-address/repository"

type WalletAddressService struct {
	IWalletAddressRepository *repository.WalletAddressRepository
}
