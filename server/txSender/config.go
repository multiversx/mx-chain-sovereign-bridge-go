package txSender

// TxSenderConfig holds tx sender config
type TxSenderConfig struct {
	BridgeSCAddress            string
	Proxy                      string
	MaxRetriesSecondsWaitNonce int
}
