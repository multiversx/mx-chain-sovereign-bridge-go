package config

import (
	"github.com/multiversx/mx-chain-sovereign-bridge-go/common"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/contracts/deploy"
)

type ContractsConfig struct {
	DeployConfig deploy.DeployConfig
	WalletConfig common.WalletConfig
}
