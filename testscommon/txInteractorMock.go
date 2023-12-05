package testscommon

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
)

// TxInteractorMock mocks TxInteractor interface
type TxInteractorMock struct {
	AddTransactionCalled          func(tx *transaction.FrontendTransaction)
	ApplySignatureCalled          func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	SendTransactionsAsBunchCalled func(ctx context.Context, bunchSize int) ([]string, error)
	IsInterfaceNilCalled          func() bool
}

// AddTransaction mocks the AddTransaction method
func (mock *TxInteractorMock) AddTransaction(tx *transaction.FrontendTransaction) {
	if mock.AddTransactionCalled != nil {
		mock.AddTransactionCalled(tx)
	}
}

// ApplySignature mocks the ApplySignature method
func (mock *TxInteractorMock) ApplySignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if mock.ApplySignatureCalled != nil {
		return mock.ApplySignatureCalled(cryptoHolder, tx)
	}
	return nil
}

// SendTransactionsAsBunch mocks the SendTransactionsAsBunch method
func (mock *TxInteractorMock) SendTransactionsAsBunch(ctx context.Context, bunchSize int) ([]string, error) {
	if mock.SendTransactionsAsBunchCalled != nil {
		return mock.SendTransactionsAsBunchCalled(ctx, bunchSize)
	}
	return nil, nil
}

// IsInterfaceNil -
func (mock *TxInteractorMock) IsInterfaceNil() bool {
	return mock == nil
}
