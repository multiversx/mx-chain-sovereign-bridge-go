package main

import (
	"fmt"

	"github.com/multiversx/mx-chain-sovereign-bridge-go/cert"
)

func main() {
	err := cert.GenerateCertFile(cert.CertificateCfg{
		CertCfg: cert.CertCfg{
			Organization: "MultiversX",
			DNSName:      "localhost",
			Availability: 10,
		},
		CertFileCfg: cert.FileCfg{
			CertFile: "certificate.crt",
			PkFile:   "private_key.pem",
		},
	})
	fmt.Println(err)
}
