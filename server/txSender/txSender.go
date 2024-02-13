package txSender

import (
	"context"
	"github.com/multiversx/mx-chain-sovereign-bridge-go/common"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	coreTx "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// TxSenderArgs holds args to create a new tx sender
type TxSenderArgs struct {
	Wallet          core.CryptoComponentsHolder
	Proxy           common.Proxy
	TxInteractor    common.TxInteractor
	TxNonceHandler  common.TxNonceSenderHandler
	DataFormatter   DataFormatter
	SCBridgeAddress string
}

type txSender struct {
	wallet          core.CryptoComponentsHolder
	netConfigs      *data.NetworkConfig
	txInteractor    common.TxInteractor
	txNonceHandler  common.TxNonceSenderHandler
	dataFormatter   DataFormatter
	scBridgeAddress string
}

// NewTxSender creates a new tx sender
func NewTxSender(args TxSenderArgs) (*txSender, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	networkConfig, err := args.Proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &txSender{
		wallet:          args.Wallet,
		netConfigs:      networkConfig,
		txInteractor:    args.TxInteractor,
		dataFormatter:   args.DataFormatter,
		scBridgeAddress: args.SCBridgeAddress,
		txNonceHandler:  args.TxNonceHandler,
	}, nil
}

func checkArgs(args TxSenderArgs) error {
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
	if len(args.SCBridgeAddress) == 0 {
		return errNoSCBridgeAddress
	}

	return nil
}

// SendTxs should send bridge data operation txs
func (ts *txSender) SendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	if len(data.Data) == 0 {
		return make([]string, 0), nil
	}

	return ts.createAndSendTxs(ctx, data)
}

func (ts *txSender) createAndSendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	txHashes := make([]string, 0)
	txsData := ts.dataFormatter.CreateTxsData(data)

	for _, txData := range txsData {
		tx := &coreTx.FrontendTransaction{
			Value:    "0",
			Receiver: ts.scBridgeAddress,
			Sender:   ts.wallet.GetBech32(),
			GasPrice: ts.netConfigs.MinGasPrice,
			GasLimit: 50_000_000, // todo
			Data:     txData,
			ChainID:  ts.netConfigs.ChainID,
			Version:  ts.netConfigs.MinTransactionVersion,
		}

		err := ts.txNonceHandler.ApplyNonceAndGasPrice(ctx, ts.wallet.GetAddressHandler(), tx)
		if err != nil {
			return nil, err
		}

		err = ts.txInteractor.ApplyUserSignature(ts.wallet, tx)
		if err != nil {
			return nil, err
		}

		hash, err := ts.txNonceHandler.SendTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}

		txHashes = append(txHashes, hash)
	}

	return txHashes, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (ts *txSender) IsInterfaceNil() bool {
	return ts == nil
}
