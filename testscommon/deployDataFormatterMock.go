package testscommon

// DeployDataFormatterMock mocks DataFormatter interface
type DeployDataFormatterMock struct {
	CreateTxsDataCalled func(data []byte) [][]byte
}

// CreateTxsData mocks the CreateTxsData method
func (mock *DeployDataFormatterMock) CreateTxsData(data []byte) [][]byte {
	if mock.CreateTxsDataCalled != nil {
		return mock.CreateTxsDataCalled(data)
	}
	return nil
}

// IsInterfaceNil -
func (mock *DeployDataFormatterMock) IsInterfaceNil() bool {
	return mock == nil
}
