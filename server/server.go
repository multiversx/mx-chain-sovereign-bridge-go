package server

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

type server struct {
	*sovereign.UnimplementedBridgeTxSenderServer
}

func NewServer() (*server, error) {
	return &server{}, nil
}

func (s *server) Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	_ = ctx
	_ = data

	return nil, nil
}

func (s *server) IsInterfaceNil() bool {
	return s == nil
}
