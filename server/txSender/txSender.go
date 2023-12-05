package txSender

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	coreTx "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	waitTimeRetrialsMs = 50
)

// TxSenderArgs holds args to create a new tx sender
type TxSenderArgs struct {
	Wallet                core.CryptoComponentsHolder
	Proxy                 Proxy
	TxInteractor          TxInteractor
	DataFormatter         DataFormatter
	SCBridgeAddress       string
	MaxRetrialsGetAccount int
}

type txSender struct {
	wallet                core.CryptoComponentsHolder
	proxy                 Proxy
	netConfigs            *data.NetworkConfig
	txInteractor          TxInteractor
	dataFormatter         DataFormatter
	scBridgeAddress       string
	waitNonce             uint64
	maxRetrialsGetAccount int
	mut                   sync.RWMutex
}

// NewTxSender creates a new tx sender
func NewTxSender(args TxSenderArgs) (*txSender, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	account, err := args.Proxy.GetAccount(context.Background(), args.Wallet.GetAddressHandler())
	if err != nil {
		return nil, err
	}

	networkConfig, err := args.Proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &txSender{
		wallet:                args.Wallet,
		proxy:                 args.Proxy,
		netConfigs:            networkConfig,
		txInteractor:          args.TxInteractor,
		dataFormatter:         args.DataFormatter,
		scBridgeAddress:       args.SCBridgeAddress,
		waitNonce:             account.Nonce,
		maxRetrialsGetAccount: args.MaxRetrialsGetAccount,
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
	if len(args.SCBridgeAddress) == 0 {
		return errNoSCBridgeAddress
	}
	if args.MaxRetrialsGetAccount == 0 {
		return errZeroTimeWaitAccountNonceUpdate
	}

	return nil
}

// SendTxs should send bridge data operation txs
func (ts *txSender) SendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	if len(data.Data) == 0 {
		return make([]string, 0), nil
	}

	ts.mut.Lock()
	defer ts.mut.Unlock()

	account, err := ts.getUpdatedAccount(ctx)
	if err != nil {
		return nil, err
	}

	numTxs, err := ts.createTxs(data, account)
	if err != nil {
		return nil, err
	}

	txHashes, err := ts.txInteractor.SendTransactionsAsBunch(ctx, numTxs)
	if err != nil {
		return nil, err
	}

	ts.waitNonce = account.Nonce + uint64(numTxs)
	return txHashes, nil
}

func (ts *txSender) createTxs(data *sovereign.BridgeOperations, account *data.Account) (int, error) {
	txsData := ts.dataFormatter.CreateTxsData(data)
	nonce := account.Nonce
	for _, txData := range txsData {
		tx := &coreTx.FrontendTransaction{
			Nonce:    nonce,
			Value:    "1",
			Receiver: ts.scBridgeAddress,
			Sender:   ts.wallet.GetBech32(),
			GasPrice: ts.netConfigs.MinGasPrice,
			GasLimit: 50_000_000, // todo
			Data:     txData,
			ChainID:  ts.netConfigs.ChainID,
			Version:  ts.netConfigs.MinTransactionVersion,
		}

		err := ts.txInteractor.ApplySignature(ts.wallet, tx)
		if err != nil {
			return 0, err
		}

		ts.txInteractor.AddTransaction(tx)
		nonce++
	}

	return int(nonce - account.Nonce), nil
}

func (ts *txSender) getUpdatedAccount(ctx context.Context) (*data.Account, error) {
	numRetrials := 0
	for numRetrials < ts.maxRetrialsGetAccount {
		acc, err := ts.proxy.GetAccount(ctx, ts.wallet.GetAddressHandler())
		if err != nil {
			log.Error("txSender.waitForNonce", "error", err)

			waitInCaseOfError(&numRetrials)
			continue
		}

		if acc.Nonce == ts.waitNonce {
			return acc, nil
		}

		log.Debug("txSender.getUpdatedAccount, waiting for account nonce update",
			"account nonce", acc.Nonce, "expected nonce", ts.waitNonce)

		time.Sleep(time.Second)
		numRetrials++
	}

	return nil, fmt.Errorf("%w after %d retrials", errCannotGetAccount, ts.maxRetrialsGetAccount)
}

func waitInCaseOfError(numRetrials *int) {
	*numRetrials++
	sleepDuration := calcRetryBackOffTime(*numRetrials)
	log.Warn("txSender.waitForNonce.proxy.GetAccount; retrying...",
		"num retrials", *numRetrials,
		"sleep duration", sleepDuration,
	)

	time.Sleep(sleepDuration)
}

func calcRetryBackOffTime(attemptNumber int) time.Duration {
	exp := math.Exp2(float64(attemptNumber))
	return time.Duration(exp) * waitTimeRetrialsMs * time.Millisecond
}

// IsInterfaceNil checks if the underlying pointer is nil
func (ts *txSender) IsInterfaceNil() bool {
	return ts == nil
}
