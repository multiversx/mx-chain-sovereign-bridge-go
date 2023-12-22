package config

import (
	"github.com/multiversx/mx-chain-sovereign-bridge-go/cert"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/server/txSender"
)

// ServerConfig holds necessary config for the grpc server
type ServerConfig struct {
	GRPCPort          string
	TxSenderConfig    txSender.TxSenderConfig
	WalletConfig      txSender.WalletConfig
	CertificateConfig cert.FileCfg
}
