package server

import (
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/cmd/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/txSender"
)

// CreateServer creates a new bridge txs sender grpc server
func CreateServer(cfg *config.ServerConfig) (sovereign.BridgeTxSenderServer, error) {
	_, err := txSender.LoadWallet(cfg.WalletConfig)
	if err != nil {
		return nil, err
	}

	return NewSovereignBridgeTxServer()
}
