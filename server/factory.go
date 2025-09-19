package server

import (
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/sovereign/dto"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/cmd/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/txSender"
)

// CreateSovereignBridgeServer creates a new bridge txs sender grpc server
func CreateSovereignBridgeServer(cfg *config.ServerConfig) (sovereign.BridgeTxSenderServer, error) {
	wallet, err := txSender.LoadWallet(cfg.WalletConfig)
	if err != nil {
		return nil, err
	}

	txSnd, err := txSender.CreateTxSender(wallet, cfg.TxSenderConfig)
	if err != nil {
		return nil, err
	}

	txSenders := map[dto.ChainID]TxSender{
		dto.MVX: txSnd,
	}

	return NewSovereignBridgeTxServer(txSenders)
}
