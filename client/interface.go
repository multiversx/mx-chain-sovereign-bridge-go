package client

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"google.golang.org/grpc"
)

// ClientHandler defines a wrapper over the grpc client connection and tx sender
type ClientHandler interface {
	Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error)
	Close() error
	IsInterfaceNil() bool
}

// GRPCConn defines a grpc client connection with closable behavior
type GRPCConn interface {
	grpc.ClientConnInterface
	Close() error
}
