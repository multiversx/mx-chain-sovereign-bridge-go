package testscommon

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
)

// TxInteractorMock mocks TxInteractor interface
type TxInteractorMock struct {
	ApplyUserSignatureCalled func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	IsInterfaceNilCalled     func() bool
}

// ApplyUserSignature -
func (mock *TxInteractorMock) ApplyUserSignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if mock.ApplyUserSignatureCalled != nil {
		return mock.ApplyUserSignatureCalled(cryptoHolder, tx)
	}
	return nil
}

// IsInterfaceNil -
func (mock *TxInteractorMock) IsInterfaceNil() bool {
	return mock == nil
}
