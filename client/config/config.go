package config

import "github.com/multiversx/mx-chain-sovereign-bridge-go/cert"

// ClientConfig holds all grpc client's config
type ClientConfig struct {
	Enabled        bool
	GRPCHost       string
	GRPCPort       string
	CertificateCfg cert.FileCfg
}
