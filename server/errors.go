package server

import "errors"

var errNilTxSender = errors.New("nil tx sender provided")

var errNilMarshaller = errors.New("nil marshaller provided")
