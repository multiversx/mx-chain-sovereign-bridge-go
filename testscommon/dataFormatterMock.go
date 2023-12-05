package testscommon

import "github.com/multiversx/mx-chain-core-go/data/sovereign"

// DataFormatterMock mocks DataFormatter interface
type DataFormatterMock struct {
	CreateTxsDataCalled  func(data *sovereign.BridgeOperations) [][]byte
	IsInterfaceNilCalled func() bool
}

// CreateTxsData mocks the CreateTxsData method
func (mock *DataFormatterMock) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	if mock.CreateTxsDataCalled != nil {
		return mock.CreateTxsDataCalled(data)
	}
	return nil
}

// IsInterfaceNil -
func (mock *DataFormatterMock) IsInterfaceNil() bool {
	return mock == nil
}
