package txSender

import (
	"context"
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	coreTx "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	gasLimitDefault       = 50_000_000
	gasLimitRegisterToken = 80_000_000
)

// TxSenderArgs holds args to create a new tx sender
type TxSenderArgs struct {
	Wallet                    core.CryptoComponentsHolder
	Proxy                     Proxy
	TxInteractor              TxInteractor
	TxNonceHandler            TxNonceSenderHandler
	DataFormatter             DataFormatter
	SCHeaderVerifierAddress   string
	SCEsdtSafeAddress         string
	SCChangeValidatorsAddress string
	SCChainConfigAddress      string
}

type txConfig struct {
	receiver string
	gasLimit uint64
}

type txSender struct {
	wallet         core.CryptoComponentsHolder
	netConfigs     *data.NetworkConfig
	txInteractor   TxInteractor
	txNonceHandler TxNonceSenderHandler
	dataFormatter  DataFormatter
	txConfigs      map[string]*txConfig
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
		wallet:         args.Wallet,
		netConfigs:     networkConfig,
		txInteractor:   args.TxInteractor,
		txNonceHandler: args.TxNonceHandler,
		dataFormatter:  args.DataFormatter,
		txConfigs: map[string]*txConfig{
			registerBridgeOpsPrefix: {
				receiver: args.SCHeaderVerifierAddress,
				gasLimit: gasLimitDefault,
			},

			executeDepositBridgeOpsPrefix: {
				receiver: args.SCEsdtSafeAddress,
				gasLimit: gasLimitDefault,
			},
			executeRegisterTokenPrefix: {
				receiver: args.SCEsdtSafeAddress,
				gasLimit: gasLimitRegisterToken,
			},

			changeValidatorSetPrefix: {
				receiver: args.SCChangeValidatorsAddress,
				gasLimit: gasLimitDefault,
			},

			executeRegisterValidatorPrefix: {
				receiver: args.SCChainConfigAddress,
				gasLimit: gasLimitDefault,
			},
			executeUnRegisterValidatorPrefix: {
				receiver: args.SCChainConfigAddress,
				gasLimit: gasLimitDefault,
			},
		},
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
	if len(args.SCHeaderVerifierAddress) == 0 {
		return errNoHeaderVerifierSCAddress
	}
	if len(args.SCEsdtSafeAddress) == 0 {
		return errNoEsdtSafeSCAddress
	}
	if len(args.SCChangeValidatorsAddress) == 0 {
		return errNoChangeValidatorSetSCAddress
	}
	if len(args.SCChainConfigAddress) == 0 {
		return errNoChainConfigSCAddress
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
			Sender:   ts.wallet.GetBech32(),
			GasPrice: ts.netConfigs.MinGasPrice,
			GasLimit: gasLimitDefault, // todo: we need proper gas estimation in the future
			Data:     txData,
			ChainID:  ts.netConfigs.ChainID,
			Version:  ts.netConfigs.MinTransactionVersion,
		}

		err := ts.setTxFields(txData, tx)
		if err != nil {
			log.Error("invalid tx data received", "data", string(txData), "error", err)
			continue
		}

		err = ts.txNonceHandler.ApplyNonceAndGasPrice(ctx, tx)
		if err != nil {
			return nil, err
		}

		err = ts.txInteractor.ApplyUserSignature(ts.wallet, tx)
		if err != nil {
			return nil, err
		}

		hash, err := ts.txNonceHandler.SendTransactions(ctx, tx)
		if err != nil {
			log.Error("failed to send tx", "error", err, "nonce", tx.Nonce)
			return nil, err
		}

		txHashes = append(txHashes, hash...)
	}

	return txHashes, nil
}

func (ts *txSender) setTxFields(txData []byte, tx *coreTx.FrontendTransaction) error {
	prefixID := getTxDataPrefix(txData)
	txCfg, found := ts.txConfigs[prefixID]
	if !found {
		return fmt.Errorf("%w, prefix = %s", errInvalidTxDataPrefix, prefixID)
	}

	tx.Receiver = txCfg.receiver
	tx.GasLimit = txCfg.gasLimit
	return nil
}

func getTxDataPrefix(txData []byte) string {
	prefix := strings.Split(string(txData), "@")
	return prefix[0]
}

// IsInterfaceNil checks if the underlying pointer is nil
func (ts *txSender) IsInterfaceNil() bool {
	return ts == nil
}
