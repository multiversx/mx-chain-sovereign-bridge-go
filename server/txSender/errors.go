package txSender

import "errors"

var errInvalidWalletType = errors.New("invalid/unknown wallet type")

var errNilWallet = errors.New("nil wallet provided")

var errNilProxy = errors.New("nil proxy provided")

var errNilTxInteractor = errors.New("nil tx interactor provided")

var errNilDataFormatter = errors.New("nil data formatter provided")

var errNilNonceHandler = errors.New("nil nonce handler provided")

var errNoSCMultisigAddress = errors.New("no sc multisig address provided")

var errNoSCEsdtSafeAddress = errors.New("no sc esdt safe address provided")
