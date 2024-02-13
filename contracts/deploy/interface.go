package deploy

// DataFormatter should format txs data for bridge operations
type DataFormatter interface {
	CreateTxsData(bytes []byte) [][]byte
	IsInterfaceNil() bool
}
