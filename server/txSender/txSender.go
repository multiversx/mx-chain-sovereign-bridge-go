package txSender

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	transaction2 "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

type TxSenderArgs struct {
	Wallet        core.CryptoComponentsHolder
	Proxy         blockchain.Proxy
	NetConfigs    *data.NetworkConfig
	TxInteractor  TxInteractor
	DataFormatter DataFormatter
	Receiver      string
}

type txSender struct {
	wallet        core.CryptoComponentsHolder
	proxy         blockchain.Proxy
	netConfigs    *data.NetworkConfig
	txInteractor  TxInteractor
	dataFormatter DataFormatter
	receiver      string
}

func NewTxSender(args TxSenderArgs) (*txSender, error) {
	return &txSender{
		wallet:        args.Wallet,
		proxy:         args.Proxy,
		netConfigs:    args.NetConfigs,
		txInteractor:  args.TxInteractor,
		dataFormatter: args.DataFormatter,
		receiver:      args.Receiver,
	}, nil
}

func (ts *txSender) SendTx(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	if len(data.Data) == 0 {
		return make([]string, 0), nil
	}

	account, err := ts.proxy.GetAccount(ctx, ts.wallet.GetAddressHandler())
	if err != nil {
		return nil, err
	}

	numTxs, err := ts.createTxs(data, account)
	if err != nil {
		return nil, err
	}

	return ts.txInteractor.SendTransactionsAsBunch(ctx, numTxs)
}

func (ts *txSender) createTxs(data *sovereign.BridgeOperations, account *data.Account) (int, error) {
	txsData := ts.dataFormatter.CreateTxsData(data)
	nonce := account.Nonce
	for _, txData := range txsData {
		tx := &transaction2.FrontendTransaction{
			Nonce:    nonce,
			Value:    "",
			Receiver: "", // todo
			Sender:   ts.wallet.GetBech32(),
			GasPrice: ts.netConfigs.MinGasPrice, // todo
			GasLimit: 300_000_000,               // todo
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

	return int(nonce-account.Nonce) + 1, nil
}
