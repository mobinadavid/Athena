package crypto

import (
	"athena/src/models"
	"fmt"
)

type ExplorerFactory struct{}

func (f *ExplorerFactory) CreateExplorer(blockchain *models.Blockchain, baseUrl string) (BlockchainExplorer, error) {
	switch blockchain.NativeAsset {
	case "ETH":
		return NewEtherscan(baseUrl)
	case "TRX":
		return NewTronscan(baseUrl)
	case "BSC":
		return NewBscscan(baseUrl)
	default:
		return nil, fmt.Errorf("unsupported blockchain: %s", blockchain.NativeAsset)
	}
}
