package server

import (
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/cmd/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/txSender"
)

// CreateServer creates a new grpc server
func CreateServer(cfg *config.ServerConfig) (sovereign.BridgeTxSenderServer, error) {
	wallet, err := txSender.LoadWallet(cfg.WalletConfig)
	if err != nil {
		return nil, err
	}

	txSnd, err := txSender.CreateTxSender(wallet, cfg.BridgeSCAddress)
	if err != nil {
		return nil, err
	}

	return NewServer(txSnd)
}
