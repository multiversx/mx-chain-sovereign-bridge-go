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

type wallet struct {
	core.CryptoComponentsHolder
}

func LoadWallet(cfg WalletConfig) (*wallet, error) {
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

	_, err = w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	holder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	if err != nil {
		return nil, err
	}

	return &wallet{
		CryptoComponentsHolder: holder,
	}, nil
}

func getWalletType(walletPath string) string {
	tokens := strings.Split(walletPath, ".")
	if len(tokens) < 2 {
		return ""
	}

	return tokens[len(tokens)-1]
}
