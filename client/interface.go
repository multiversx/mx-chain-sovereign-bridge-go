package client

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"google.golang.org/grpc"
)

type ClientHandler interface {
	Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error)
	Close() error
	IsInterfaceNil() bool
}

type GRPCConn interface {
	grpc.ClientConnInterface
	Close() error
}
