package txSender

import (
	"context"

	transaction2 "github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

type TxSenderArgs struct {
	Wallet       core.CryptoComponentsHolder
	Proxy        blockchain.Proxy
	NetConfigs   *data.NetworkConfig
	TxInteractor TxInteractor
}

type txSender struct {
	wallet       core.CryptoComponentsHolder
	proxy        blockchain.Proxy
	netConfigs   *data.NetworkConfig
	txInteractor TxInteractor
}

func NewTxSender(args TxSenderArgs) (*txSender, error) {
	return &txSender{
		wallet:       args.Wallet,
		proxy:        args.Proxy,
		netConfigs:   args.NetConfigs,
		txInteractor: args.TxInteractor,
	}, nil
}

func (ts *txSender) SendTx() error {
	account, err := ts.proxy.GetAccount(context.Background(), ts.wallet.GetAddressHandler())
	if err != nil {
		return err
	}

	tx := &transaction2.FrontendTransaction{
		Nonce:     account.Nonce,
		Value:     "", // todo
		Receiver:  "", // todo
		Sender:    ts.wallet.GetBech32(),
		GasPrice:  0,   // todo
		GasLimit:  0,   // todo
		Data:      nil, // todo
		Signature: "",
		ChainID:   ts.netConfigs.ChainID,
		Version:   ts.netConfigs.MinTransactionVersion,
		Options:   0,
	}

	err = ts.txInteractor.ApplySignature(ts.wallet, tx)
	if err != nil {
		return err
	}
	ts.txInteractor.AddTransaction(tx)

	hashes, err := ts.txInteractor.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		return err
	}

	_ = hashes
	return nil
}
