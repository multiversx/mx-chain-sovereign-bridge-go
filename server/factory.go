package server

import (
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/common"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/cmd/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/txSender"
)

// CreateSovereignBridgeServer creates a new bridge txs sender grpc server
func CreateSovereignBridgeServer(cfg *config.ServerConfig) (sovereign.BridgeTxSenderServer, error) {
	wallet, err := common.LoadWallet(cfg.WalletConfig)
	if err != nil {
		return nil, err
	}

	txSnd, err := txSender.CreateTxSender(wallet, cfg.TxSenderConfig)
	if err != nil {
		return nil, err
	}

	return NewSovereignBridgeTxServer(txSnd)
}
