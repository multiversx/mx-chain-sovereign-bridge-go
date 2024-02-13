package deploy

type dataFormatter struct {
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter() *dataFormatter {
	return &dataFormatter{}
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(bytes []byte) [][]byte {
	txsData := make([][]byte, 0)
	if len(bytes) == 0 {
		return txsData
	}

	txsData = append(txsData, bytes)
	txsData = append(txsData, VmTypeWasmVm)
	txsData = append(txsData, CodeMetadata)

	return txsData
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}
