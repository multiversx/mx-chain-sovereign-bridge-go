package config

import "github.com/multiversx/mx-chain-sovereign-bridge-go/bridgeContracts/deploy"

type BridgeConfig struct {
	DeployConfig deploy.DeployConfig
	WalletConfig deploy.WalletConfig
}
