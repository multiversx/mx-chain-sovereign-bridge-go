package config

import "google.golang.org/grpc"

type ClientConfig struct {
	GRPCHost string
	GRPCPort string
	GRPConn  grpc.ClientConnInterface
}
