package auth

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type Certificate struct {
	CAFile     string
	CertFile   string
	KeyFile    string
	ServerName string
}

func (c Certificate) GetTLSConfigForServer() (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(
		c.CertFile,
		c.KeyFile,
	)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(c.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, errors.New("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	return tlsConfig, nil
}

func (c Certificate) GetTransportCredsForClient() (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(
		c.CertFile,
		c.KeyFile,
	)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(c.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, errors.New("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   c.ServerName,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	return transportCreds, nil
}
