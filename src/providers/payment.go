package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services"

	"github.com/google/wire"
)

var PaymentContainer = wire.NewSet(
	ProvidePaymentRequestRepository,
	ProvidePaymentService,
	ProvidePaymentController,
	wire.Bind(new(repositories.IPaymentRequestRepository), new(*repositories.PaymentRequestRepository)),
	wire.Bind(new(services.IPaymentService), new(*services.PaymentService)),
)

func ProvidePaymentRequestRepository(db *database.Database) *repositories.PaymentRequestRepository {
	return &repositories.PaymentRequestRepository{
		DatabaseHandler: db,
	}
}

func ProvidePaymentService(
	paymentRequestRepository repositories.IPaymentRequestRepository,
	walletAddressService services.IWalletAddressService,
	walletAddressRepository repositories.IWalletAddressRepository,
	blockchainService services.IBlockChainService,
	depositRepository repositories.IDepositRepository,
) *services.PaymentService {
	return &services.PaymentService{
		PaymentRequestRepository: paymentRequestRepository,
		WalletAddressService:     walletAddressService,
		WalletAddressRepository:  walletAddressRepository,
		BlockchainService:        blockchainService,
		DepositRepository:        depositRepository,
	}
}

func ProvidePaymentController(service services.IPaymentService) *controllers.PaymentController {
	return &controllers.PaymentController{
		PaymentService: service,
	}
}
