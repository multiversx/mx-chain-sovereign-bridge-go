package main

import (
	"fmt"

	"github.com/multiversx/mx-chain-sovereign-bridge-go/cert"
)

func main() {
	err := cert.GenerateCertFile()
	fmt.Println(err)
}
