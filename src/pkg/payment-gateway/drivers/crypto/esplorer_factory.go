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

func (f *ExplorerFactory) CreateExplorer(blockchain *models.Blockchain, baseUrl string) (BlockchainExplorer, error) {
	switch blockchain.NativeAsset {
	case "ETH":
		return etherscan.NewEtherscan(baseUrl)
	case "TRX":
		return tronscan.NewTronscan(baseUrl)
	case "BSC":
		return bscscan.NewBscscan(baseUrl)
	case "BTC":
		return btcscan.NewBtcscan(baseUrl)
	default:
		return nil, fmt.Errorf("unsupported blockchain: %s", blockchain.NativeAsset)
	}
}
