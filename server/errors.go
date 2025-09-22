package server

import "errors"

var errNilTxSenders = errors.New("nil tx senders provided")

var errNilTxSender = errors.New("nil tx sender provided")

var errNilMarshaller = errors.New("nil marshaller provided")

var errNilGinHandler = errors.New("nil gin handler provided")

var errNilGRPCHandler = errors.New("nil grpc handler provided")
