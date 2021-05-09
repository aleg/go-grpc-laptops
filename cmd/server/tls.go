package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

const (
	serverCertFile   = "cert/server-cert.pem"
	serverKeyFile    = "cert/server-key.pem"
	clientCACertFile = "cert/ca-cert.pem"
)

func loadTlsCred() (credentials.TransportCredentials, error) {
	// gRPC connection types:
	// - server-side TLS: only the server needs to provide its
	//   certificate to the client;
	// - mutual TLS: both server and client need to provide
	//   certificates to each other.

	// -------------------------------------------------------------
	// Needed for mutual TLS.
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(clientCACertFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("Failed to add client CA's certificate")
	}

	// -------------------------------------------------------------
	// Needed for server-side TLS.
	// Load server's certificate and private key.
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------
	// Create the transport credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		// ClientAuth:   tls.NoClientCert,// server-side only TLS
		ClientAuth: tls.RequireAndVerifyClientCert, // mutual TLS
		ClientCAs:  certPool,                       // mutual TLS
	}

	return credentials.NewTLS(config), nil
}
