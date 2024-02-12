package deploy

// WalletConfig holds wallet config
type WalletConfig struct {
	Path     string
	Password string
}

// DeployConfig holds deploy config
type DeployConfig struct {
	Proxy                      string
	MaxRetriesSecondsWaitNonce int
}
