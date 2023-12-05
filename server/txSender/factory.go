package txSender

import (
	"time"

	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

func CreateTxSender(wallet core.CryptoComponentsHolder, cfg TxSenderConfig) (*txSender, error) {
	args := blockchain.ArgsProxy{
		ProxyURL:            cfg.Proxy,
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

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		return nil, err
	}

	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		return nil, err
	}

	return NewTxSender(TxSenderArgs{
		Wallet:                wallet,
		Proxy:                 proxy,
		TxInteractor:          ti,
		DataFormatter:         NewDataFormatter(),
		SCBridgeAddress:       cfg.BridgeSCAddress,
		MaxRetrialsGetAccount: cfg.MaxRetrialsWaitNonce,
	})
}
