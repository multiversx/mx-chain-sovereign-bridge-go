package common

import "errors"

var errInvalidWalletType = errors.New("invalid/unknown wallet type")

var ErrNilWallet = errors.New("nil wallet provided")

var ErrNilProxy = errors.New("nil proxy provided")

var ErrNilTxInteractor = errors.New("nil tx interactor provided")

var ErrNilDataFormatter = errors.New("nil data formatter provided")

var ErrNilNonceHandler = errors.New("nil nonce handler provided")
