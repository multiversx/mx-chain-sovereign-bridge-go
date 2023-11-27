package txSender

import (
	"context"
	"time"

	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
)

func CreateTxSender(wallet core.CryptoComponentsHolder) (*txSender, error) {
	args := blockchain.ArgsProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	proxy, err := blockchain.NewProxy(args)
	if err != nil {
		return nil, err
	}

	netConfigs, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		return nil, err
	}

	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		return nil, err
	}

	return NewTxSender(TxSenderArgs{
		Wallet:       wallet,
		Proxy:        proxy,
		NetConfigs:   netConfigs,
		TxInteractor: ti,
	})
}
