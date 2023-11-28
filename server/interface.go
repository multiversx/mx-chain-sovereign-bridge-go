package server

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

// TxSender defines a tx sender for bridge operations
type TxSender interface {
	SendTx(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error)
	IsInterfaceNil() bool
}
