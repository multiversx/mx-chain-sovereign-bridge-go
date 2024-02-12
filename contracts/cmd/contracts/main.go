package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/contracts/cmd/config"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/contracts/deploy"
	"github.com/urfave/cli"
)

var log = logger.GetOrCreate("sov-bridge-contracts")

const (
	envWallet              = "WALLET_PATH"
	envPassword            = "WALLET_PASSWORD"
	envMultiversXProxy     = "MULTIVERSX_PROXY"
	envMaxRetriesWaitNonce = "MAX_RETRIES_SECONDS_WAIT_NONCE"
)

func main() {
	app := cli.NewApp()
	app.Name = "Sovereign bridge contracts deploy"
	app.Usage = ""
	app.Action = deployBridgeContracts
	app.Flags = []cli.Flag{}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func deployBridgeContracts(ctx *cli.Context) error {
	_, err := loadConfig()
	if err != nil {
		return err
	}

	log.Info("sovereign bridge contracts deployed successfully")

	return nil
}

func loadConfig() (*config.ContractsConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	walletPath := os.Getenv(envWallet)
	walletPassword := os.Getenv(envPassword)
	proxy := os.Getenv(envMultiversXProxy)
	maxRetriesWaitNonceStr := os.Getenv(envMaxRetriesWaitNonce)

	maxRetriesWaitNonce, err := strconv.Atoi(maxRetriesWaitNonceStr)
	if err != nil {
		return nil, err
	}

	log.Info("loaded config", "proxy", proxy)
	log.Info("loaded config", "maxRetriesWaitNonce", maxRetriesWaitNonce)

	return &config.ContractsConfig{
		WalletConfig: deploy.WalletConfig{
			Path:     walletPath,
			Password: walletPassword,
		},
		DeployConfig: deploy.DeployConfig{
			Proxy:                      proxy,
			MaxRetriesSecondsWaitNonce: maxRetriesWaitNonce,
		},
	}, nil
}
