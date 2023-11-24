package client

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
)

type client struct {
	bridgeClient sovereign.BridgeTxSenderClient
	conn         GRPCConn
}

func NewClient(bridgeClient sovereign.BridgeTxSenderClient, conn GRPCConn) (*client, error) {
	if conn == nil {
		return nil, errNilClientConnection
	}
	if bridgeClient == nil {
		return nil, nil
	}

	return &client{
		conn:         conn,
		bridgeClient: bridgeClient,
	}, nil
}

func (c *client) Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	return c.bridgeClient.Send(ctx, data)
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) IsInterfaceNil() bool {
	return c == nil
}
