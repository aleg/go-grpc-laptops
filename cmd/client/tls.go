package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

const (
	serverCAFile = "cert/ca-cert.pem"
)

func loadTlsCred() (credentials.TransportCredentials, error) {
	// -------------------------------------------------------------
	// Needed for server-side TLS.
	// Load certificate of the CA who signed server's certificate.
	pemServerCA, err := ioutil.ReadFile(serverCAFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("Failed to add server CA's certificate")
	}

	// -------------------------------------------------------------
	// Needed for mutual TLS.
	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------
	// Create the transport credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert}, // mutual TLS
		RootCAs:      certPool,                      // server-side TLS
	}
	return credentials.NewTLS(config), nil
}
