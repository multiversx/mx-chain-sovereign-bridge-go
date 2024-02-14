package deploy

import (
	"context"
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreTx "github.com/multiversx/mx-chain-core-go/data/transaction"
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

func (d *deployerArgs) DeployEsdtSafeContract(ctx context.Context, contractLocation string) error {
	log.Info("deploying esdt-safe contract")

	binary, err := readContractWasm(contractLocation)
	if err != nil {
		return err
	}

	log.Info("esdt-safe contract", "size", len(binary))

	hash, err := d.sendDeployTx(ctx, binary)
	if err != nil {
		return err
	}

	log.Info("esdt-safe contract deployed", "hash", hash)

	return nil
}

func (d *deployerArgs) sendDeployTx(ctx context.Context, wasmBinary []byte) (string, error) {
	if len(wasmBinary) == 0 {
		return "", nil
	}

	return d.createAndSendTx(ctx, wasmBinary)
}

func (d *deployerArgs) createAndSendTx(ctx context.Context, wasmBinary []byte) (string, error) {
	txData := d.dataFormatter.CreateTxsData(wasmBinary)

	tx := &coreTx.FrontendTransaction{
		Value:    "0",
		Receiver: SystemScAddress,
		Sender:   d.wallet.GetBech32(),
		GasPrice: d.netConfigs.MinGasPrice,
		GasLimit: 50_000_000,
		Data:     txData,
		ChainID:  d.netConfigs.ChainID,
		Version:  2,
	}

	err := d.txNonceHandler.ApplyNonceAndGasPrice(ctx, d.wallet.GetAddressHandler(), tx)
	if err != nil {
		return "", err
	}

	err = d.txInteractor.ApplyUserSignature(d.wallet, tx)
	if err != nil {
		return "", err
	}

	hash, err := d.txNonceHandler.SendTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (d *deployerArgs) IsInterfaceNil() bool {
	return d == nil
}
