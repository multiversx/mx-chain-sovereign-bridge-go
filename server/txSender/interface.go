package txSender

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/sovereign"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// TxInteractor defines a tx interactor with multiversx blockchain
type TxInteractor interface {
	AddTransaction(tx *transaction.FrontendTransaction)
	ApplySignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	SendTransactionsAsBunch(ctx context.Context, bunchSize int) ([]string, error)
	IsInterfaceNil() bool
}

// Proxy defines the proxy to interact with MultiversX blockchain
type Proxy interface {
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	IsInterfaceNil() bool
}

// DataFormatter should format txs data for bridge operations
type DataFormatter interface {
	CreateTxsData(data *sovereign.BridgeOperations) [][]byte
	IsInterfaceNil() bool
}

type TxNonceSenderHandler interface {
	GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error)
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	IsInterfaceNil() bool
}
