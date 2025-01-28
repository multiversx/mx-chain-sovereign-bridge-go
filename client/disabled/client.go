package disabled

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

type client struct{}

// NewClient creates a new instance of disabled client
func NewClient() *client {
	return &client{}
}

// Send -
func (c *client) Send(_ context.Context, _ *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	return &sovereign.BridgeOperationsResponse{}, nil
}

// Close -
func (c *client) Close() error {
	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (c *client) IsInterfaceNil() bool {
	return c == nil
}
