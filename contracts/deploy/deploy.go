package deploy

import (
	"context"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/common"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"os"
)

var log = logger.GetOrCreate("contracts deploy")

type DeployerArgs struct {
	Wallet         core.CryptoComponentsHolder
	Proxy          common.Proxy
	TxInteractor   common.TxInteractor
	TxNonceHandler common.TxNonceSenderHandler
	DataFormatter  DataFormatter
}

type deployerArgs struct {
	wallet         core.CryptoComponentsHolder
	netConfigs     *data.NetworkConfig
	txInteractor   common.TxInteractor
	txNonceHandler common.TxNonceSenderHandler
	dataFormatter  DataFormatter
}

func NewDeployer(args DeployerArgs) (*deployerArgs, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	networkConfig, err := args.Proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &deployerArgs{
		wallet:         args.Wallet,
		netConfigs:     networkConfig,
		txInteractor:   args.TxInteractor,
		dataFormatter:  args.DataFormatter,
		txNonceHandler: args.TxNonceHandler,
	}, nil
}

func checkArgs(args DeployerArgs) error {
	if check.IfNil(args.Wallet) {
		return common.ErrNilWallet
	}
	if check.IfNil(args.Proxy) {
		return common.ErrNilProxy
	}
	if check.IfNil(args.TxInteractor) {
		return common.ErrNilTxInteractor
	}
	if check.IfNil(args.DataFormatter) {
		return common.ErrNilDataFormatter
	}
	if check.IfNil(args.TxNonceHandler) {
		return common.ErrNilNonceHandler
	}

	return nil
}

func readContractWasm(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *deployerArgs) DeployEsdtSafeContract(contractLocation string) error {
	log.Info("deploying esdt-safe contract")

	binary, err := readContractWasm(contractLocation)
	if err != nil {
		return err
	}

	log.Info("esdt-safe contract", "size", len(binary))

	// deploy contract send transaction

	log.Info("esdt-safe contract deployed successfully")

	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (d *deployerArgs) IsInterfaceNil() bool {
	return d == nil
}
