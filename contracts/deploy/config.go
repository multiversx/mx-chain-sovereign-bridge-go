package deploy

// DeployConfig holds deploy config
type DeployConfig struct {
	Contracts                  ContractsLocation
	Proxy                      string
	MaxRetriesSecondsWaitNonce int
}

// ContractsLocation holds all contracts path
type ContractsLocation struct {
	EsdtSafeContractPath string
}
