package txSender

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	transaction2 "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	maxRequestsRetrial = 10
	waitTimeRetrialsMs = 50
)

// TxSenderArgs holds args to create a new tx sender
type TxSenderArgs struct {
	Wallet          core.CryptoComponentsHolder
	Proxy           blockchain.Proxy
	NetConfigs      *data.NetworkConfig
	TxInteractor    TxInteractor
	DataFormatter   DataFormatter
	SCBridgeAddress string
}

type txSender struct {
	wallet          core.CryptoComponentsHolder
	proxy           blockchain.Proxy
	netConfigs      *data.NetworkConfig
	txInteractor    TxInteractor
	dataFormatter   DataFormatter
	scBridgeAddress string
	waitNonce       uint64
	mut             sync.RWMutex
}

// NewTxSender creates a new tx sender
func NewTxSender(args TxSenderArgs) (*txSender, error) {
	if check.IfNil(args.Wallet) {
		return nil, errNilWallet
	}
	if check.IfNil(args.Proxy) {
		return nil, errNilProxy
	}
	if check.IfNil(args.TxInteractor) {
		return nil, errNilTxInteractor
	}
	if check.IfNil(args.DataFormatter) {
		return nil, errNilDataFormatter
	}
	if len(args.SCBridgeAddress) == 0 {
		return nil, errNoSCBridgeAddress
	}
	if args.NetConfigs == nil {
		return nil, errNilNetworkConfigs
	}

	account, err := args.Proxy.GetAccount(context.Background(), args.Wallet.GetAddressHandler())
	if err != nil {
		return nil, err
	}

	ts := &txSender{
		wallet:          args.Wallet,
		proxy:           args.Proxy,
		netConfigs:      args.NetConfigs,
		txInteractor:    args.TxInteractor,
		dataFormatter:   args.DataFormatter,
		scBridgeAddress: args.SCBridgeAddress,
		waitNonce:       account.Nonce,
	}

	return ts, nil
}

// SendTx should send bridge data operation txs
func (ts *txSender) SendTx(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
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

	ts.waitNonce = account.Nonce + uint64(numTxs)
	return ts.txInteractor.SendTransactionsAsBunch(ctx, numTxs)
}

func (ts *txSender) createTxs(data *sovereign.BridgeOperations, account *data.Account) (int, error) {
	txsData := ts.dataFormatter.CreateTxsData(data)
	nonce := account.Nonce
	for _, txData := range txsData {
		tx := &transaction2.FrontendTransaction{
			Nonce:    nonce,
			Value:    "1",
			Receiver: ts.scBridgeAddress, // todo
			Sender:   ts.wallet.GetBech32(),
			GasPrice: ts.netConfigs.MinGasPrice, // todo
			GasLimit: 50_000_000,                // todo
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
	for numRetrials < maxRequestsRetrial {
		acc, err := ts.proxy.GetAccount(ctx, ts.wallet.GetAddressHandler())
		if err != nil {
			log.Error("txSender.waitForNonce", "error", err)

			waitInCaseOfError(&numRetrials)
			continue
		}

		if acc.Nonce == ts.waitNonce {
			return acc, nil
		}

		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("%w after %d retrials", errCannotGetAccount, maxRequestsRetrial)
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