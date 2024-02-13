package contracts

import (
	"context"

	"github.com/multiversx/mx-chain-sovereign-bridge-go/common"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/contracts/cmd/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/contracts/deploy"
)

func DeployBridgeContracts(ctx context.Context, cfg *config.ContractsConfig) error {
	wallet, err := common.LoadWallet(cfg.WalletConfig)
	if err != nil {
		return err
	}

	deployer, err := deploy.CreateDeployer(wallet, cfg.DeployConfig)
	if err != nil {
		return err
	}

	err = deployer.DeployEsdtSafeContract(ctx, cfg.DeployConfig.Contracts.EsdtSafeContractPath)
	if err != nil {
		return err
	}

	return nil
}