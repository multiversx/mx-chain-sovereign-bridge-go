package server

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

type TxSender interface {
	SendTx(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error)
}
