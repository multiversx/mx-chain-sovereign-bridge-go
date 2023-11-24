package client

import "errors"

var errNilClientConnection = errors.New("nil grpc client connection provided")

var errCannotOpenConnection = errors.New("cannot open connection")
