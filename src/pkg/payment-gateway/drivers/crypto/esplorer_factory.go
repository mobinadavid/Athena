package crypto

import (
	"athena/src/models"
	"athena/src/pkg/payment-gateway/drivers/crypto/bscscan"
	"athena/src/pkg/payment-gateway/drivers/crypto/btcscan"
	"athena/src/pkg/payment-gateway/drivers/crypto/etherscan"
	"athena/src/pkg/payment-gateway/drivers/crypto/tronscan"
	"fmt"
)

type ExplorerFactory struct{}

func (f *ExplorerFactory) CreateExplorer(blockchain *models.Blockchain, baseUrl string, page, limit uint) (BlockchainExplorer, error) {
	switch blockchain.NativeAsset {
	case "ETH":
		return etherscan.NewEtherscan(baseUrl, page, limit)
	case "TRX":
		return tronscan.NewTronscan(baseUrl, page, limit)
	case "BSC":
		return bscscan.NewBscscan(baseUrl, page, limit)
	case "BTC":
		return btcscan.NewBtcscan(baseUrl, page, limit)
	default:
		return nil, fmt.Errorf("unsupported blockchain: %s", blockchain.NativeAsset)
	}
}
