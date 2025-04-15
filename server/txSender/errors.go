package txSender

import "errors"

var errInvalidWalletType = errors.New("invalid/unknown wallet type")

var errNilWallet = errors.New("nil wallet provided")

var errNilProxy = errors.New("nil proxy provided")

var errNilTxInteractor = errors.New("nil tx interactor provided")

var errNilDataFormatter = errors.New("nil data formatter provided")

var errNilNonceHandler = errors.New("nil nonce handler provided")

var errNoHeaderVerifierSCAddress = errors.New("no header verifier sc address provided")

var errNoEsdtSafeSCAddress = errors.New("no esdt safe sc address provided")

var errNoChangeValidatorSetSCAddress = errors.New("no change validator set sc address provided")

var errInvalidBridgeDataSetValidatorChange = errors.New("invalid number of bridge data operations for validator set change")

var errInvalidTxDataPrefix = errors.New("invalid/unknown tx data endpoint to call")
