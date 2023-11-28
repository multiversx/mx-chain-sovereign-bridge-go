package txSender

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
)

type TxInteractor interface {
	AddTransaction(tx *transaction.FrontendTransaction)
	ApplySignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	SendTransactionsAsBunch(ctx context.Context, bunchSize int) ([]string, error)
}

type DataFormatter interface {
	CreateTxsData(data *sovereign.BridgeOperations) [][]byte
}
