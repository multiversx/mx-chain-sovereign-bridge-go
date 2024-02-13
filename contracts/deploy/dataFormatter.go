package deploy

import "encoding/hex"

type dataFormatter struct {
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter() *dataFormatter {
	return &dataFormatter{}
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(bytes []byte) []byte {
	if len(bytes) == 0 {
		return make([]byte, 0)
	}

	txData := []byte(
		hex.EncodeToString(bytes) + "@" +
			hex.EncodeToString(VmTypeWasmVm) + "@" +
			hex.EncodeToString(CodeMetadata))

	return txData
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}
