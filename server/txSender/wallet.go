package txSender

import (
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
	log    = logger.GetOrCreate("mx-chain-sovereign-bridge-go")
)

const (
	json = "json"
	pem  = "pem"
)

func LoadWallet(cfg WalletConfig) (core.CryptoComponentsHolder, error) {
	var privateKey []byte
	var err error

	w := interactors.NewWallet()
	walletType := getWalletType(cfg.Path)
	switch walletType {
	case pem:
		privateKey, err = w.LoadPrivateKeyFromPemFile(cfg.Path)
	case json:
		privateKey, err = w.LoadPrivateKeyFromJsonFile(cfg.Path, cfg.Password)
	default:
		return nil, fmt.Errorf("%w: %s, acceptable:%s, %s", errInvalidWallet, walletType, pem, json)
	}

	if err != nil {
		return nil, err
	}

	_, err = w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
}

func getWalletType(walletPath string) string {
	tokens := strings.Split(walletPath, ".")
	if len(tokens) < 2 {
		return ""
	}

	return tokens[len(tokens)-1]
}
