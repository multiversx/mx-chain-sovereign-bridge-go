package txSender

import "errors"

var errInvalidWalletType = errors.New("invalid/unknown wallet type")

var errNilWallet = errors.New("nil wallet provided")

var errNilProxy = errors.New("nil proxy provided")

var errNilTxInteractor = errors.New("nil tx interactor provided")

var errNilDataFormatter = errors.New("nil data formatter provided")

var errNoSCBridgeAddress = errors.New("no sc bridge address provided")

var errNilNetworkConfigs = errors.New("nil network configs provided")

var errCannotGetAccount = errors.New("could not get account from proxy")
