package txSender

import (
	"context"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	coreTx "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// TxSenderArgs holds args to create a new tx sender
type TxSenderArgs struct {
	Wallet            core.CryptoComponentsHolder
	Proxy             Proxy
	TxInteractor      TxInteractor
	TxNonceHandler    TxNonceSenderHandler
	DataFormatter     DataFormatter
	SCMultiSigAddress string
	SCEsdtSafeAddress string
}

type txSender struct {
	wallet            core.CryptoComponentsHolder
	netConfigs        *data.NetworkConfig
	txInteractor      TxInteractor
	txNonceHandler    TxNonceSenderHandler
	dataFormatter     DataFormatter
	scMultisigAddress string
	scEsdtSafeAddress string
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
		wallet:            args.Wallet,
		netConfigs:        networkConfig,
		txInteractor:      args.TxInteractor,
		txNonceHandler:    args.TxNonceHandler,
		dataFormatter:     args.DataFormatter,
		scMultisigAddress: args.SCMultiSigAddress,
		scEsdtSafeAddress: args.SCEsdtSafeAddress,
	}, nil
}

func checkArgs(args TxSenderArgs) error {
	if check.IfNil(args.Wallet) {
		return errNilWallet
	}
	if check.IfNil(args.Proxy) {
		return errNilProxy
	}
	if check.IfNil(args.TxInteractor) {
		return errNilTxInteractor
	}
	if check.IfNil(args.DataFormatter) {
		return errNilDataFormatter
	}
	if check.IfNil(args.TxNonceHandler) {
		return errNilNonceHandler
	}
	if len(args.SCMultiSigAddress) == 0 {
		return errNoMultiSigSCAddress
	}
	if len(args.SCEsdtSafeAddress) == 0 {
		return errNoEsdtSafeSCAddress
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
		var tx *coreTx.FrontendTransaction

		switch {
		case strings.HasPrefix(string(txData), registerBridgeOpsPrefix):
			tx = &coreTx.FrontendTransaction{
				Value:    "0",
				Receiver: ts.scMultisigAddress,
				Sender:   ts.wallet.GetBech32(),
				GasPrice: ts.netConfigs.MinGasPrice,
				GasLimit: 50_000_000, // todo
				Data:     txData,
				ChainID:  ts.netConfigs.ChainID,
				Version:  ts.netConfigs.MinTransactionVersion,
			}
		case strings.HasPrefix(string(txData), executeBridgeOpsPrefix):
			tx = &coreTx.FrontendTransaction{
				Value:    "0",
				Receiver: ts.scEsdtSafeAddress,
				Sender:   ts.wallet.GetBech32(),
				GasPrice: ts.netConfigs.MinGasPrice,
				GasLimit: 50_000_000, // todo
				Data:     txData,
				ChainID:  ts.netConfigs.ChainID,
				Version:  ts.netConfigs.MinTransactionVersion,
			}
		default:
			log.Error("invalid tx data received", "data", string(tx.Data))
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
