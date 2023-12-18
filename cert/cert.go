package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("cert")

type CertificateCfg struct {
	CertCfg     CertCfg
	CertFileCfg FileCfg
}

type CertCfg struct {
	Organization string
	DNSName      string
	Availability int64
}

type FileCfg struct {
	CertFile string
	PkFile   string
}

func GenerateCert(cfg CertCfg) ([]byte, *rsa.PrivateKey, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{cfg.Organization},
			CommonName:   cfg.Organization,
		},
		DNSNames:              []string{cfg.DNSName},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(cfg.Availability) * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, pk.Public(), pk)
	if err != nil {
		return nil, nil, err
	}

	return cert, pk, nil
}

func GenerateCertFile(cfg CertificateCfg) error {
	cert, pk, err := GenerateCert(cfg.CertCfg)
	if err != nil {
		return err
	}

	certOut, err := os.Create(cfg.CertFileCfg.CertFile)
	if err != nil {
		return err
	}
	defer func() {
		err = certOut.Close()
		log.LogIfError(err)
	}()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		return err
	}

	keyOut, err := os.Create(cfg.CertFileCfg.PkFile)
	if err != nil {
		return err
	}
	defer func() {
		err = keyOut.Close()
		log.LogIfError(err)
	}()

	pkBytes := x509.MarshalPKCS1PrivateKey(pk)
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: pkBytes})
	if err != nil {
		return err
	}

	return nil
}

func CreateTLSServerConfig(cfg FileCfg) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.PkFile)
	if err != nil {
		return nil, err
	}

	certPool, err := createCertPool(cert)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}

func CreateTLSClientConfig(cfg FileCfg) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.PkFile)
	if err != nil {
		return nil, err
	}

	certPool, err := createCertPool(cert)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}, nil
}

func createCertPool(cert tls.Certificate) (*x509.CertPool, error) {
	certLeaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certLeaf)

	return certPool, nil
}
