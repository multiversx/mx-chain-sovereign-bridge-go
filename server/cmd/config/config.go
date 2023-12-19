package config

import "github.com/multiversx/mx-chain-sovereign-bridge-go/server/txSender"

// ServerConfig holds necessary config for the grpc server
type ServerConfig struct {
	GRPCPort     string
	WalletConfig txSender.WalletConfig
}
