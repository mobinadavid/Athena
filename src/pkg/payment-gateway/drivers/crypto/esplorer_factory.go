package crypto

import (
	"athena/src/models"
	"fmt"
)

type ExplorerFactory struct{}

func (f *ExplorerFactory) CreateExplorer(blockchain *models.Blockchain, baseUrl string, page, limit uint) (BlockchainExplorer, error) {
	switch blockchain.NativeAsset {
	case "ETH":
		return NewEtherscan(baseUrl, page, limit)
	case "TRX":
		return NewTronscan(baseUrl, page, limit)
	case "BSC":
		return NewBscscan(baseUrl, page, limit)
	default:
		return nil, fmt.Errorf("unsupported blockchain: %s", blockchain.NativeAsset)
	}
}
