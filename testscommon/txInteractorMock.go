package testscommon

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
)

// TxInteractorMock mocks TxInteractor interface
type TxInteractorMock struct {
	ApplySignatureCalled func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	IsInterfaceNilCalled func() bool
}

// ApplySignature mocks the ApplySignature method
func (mock *TxInteractorMock) ApplySignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if mock.ApplySignatureCalled != nil {
		return mock.ApplySignatureCalled(cryptoHolder, tx)
	}
	return nil
}

// IsInterfaceNil -
func (mock *TxInteractorMock) IsInterfaceNil() bool {
	return mock == nil
}
