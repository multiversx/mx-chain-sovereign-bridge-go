package deploy

import (
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-sdk-go/interactors/nonceHandlerV2"
	"time"
)

func CreateDeployer(wallet core.CryptoComponentsHolder, cfg DeployConfig) (*deployerArgs, error) {
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

	nonceHandler, err := nonceHandlerV2.NewNonceTransactionHandlerV2(nonceHandlerV2.ArgsNonceTransactionsHandlerV2{
		Proxy:            proxy,
		IntervalToResend: time.Second * time.Duration(cfg.MaxRetriesSecondsWaitNonce),
	})
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

	return NewDeployer(DeployerArgs{
		Wallet:         wallet,
		Proxy:          proxy,
		TxInteractor:   ti,
		TxNonceHandler: nonceHandler,
		DataFormatter:  NewDataFormatter(),
	})
}
