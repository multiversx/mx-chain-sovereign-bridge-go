package server

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

type server struct {
	*sovereign.UnimplementedBridgeTxSenderServer
}

// NewSovereignBridgeTxServer creates a new sovereign bridge operations server. This server receives bridge data operations from
// sovereign nodes and sends transactions to main chain.
func NewSovereignBridgeTxServer() (*server, error) {
	return &server{}, nil
}

// Send should handle receiving data bridge operations from sovereign shard and forward transactions to main chain
func (s *server) Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	_ = ctx
	_ = data

	return nil, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (s *server) IsInterfaceNil() bool {
	return s == nil
}
