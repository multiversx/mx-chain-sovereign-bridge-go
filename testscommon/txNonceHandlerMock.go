package testscommon

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
)

// TxNonceSenderHandlerMock mocks TxNonceSenderHandler interface
type TxNonceSenderHandlerMock struct {
	ApplyNonceAndGasPriceCalled func(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error
	SendTransactionCalled       func(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
}

// ApplyNonceAndGasPrice mocks the ApplyNonceAndGasPrice method
func (mock *TxNonceSenderHandlerMock) ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error {
	if mock.ApplyNonceAndGasPriceCalled != nil {
		return mock.ApplyNonceAndGasPriceCalled(ctx, address, tx)
	}
	return nil
}

// SendTransaction mocks the SendTransaction method
func (mock *TxNonceSenderHandlerMock) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if mock.SendTransactionCalled != nil {
		return mock.SendTransactionCalled(ctx, tx)
	}
	return "", nil
}

// IsInterfaceNil -
func (mock *TxNonceSenderHandlerMock) IsInterfaceNil() bool {
	return mock == nil
}
