package txSender

import (
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

// DataFormatter should format txs data for bridge operations
type DataFormatter interface {
	CreateTxsData(data *sovereign.BridgeOperations) [][]byte
	IsInterfaceNil() bool
}
