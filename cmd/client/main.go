package main

import (
	"flag"
	"log"
	"time"

	"github.com/aleg/go-grpc-laptops/client"
	"google.golang.org/grpc"
)

const (
	username        = "jay"
	password        = "secret-jay"
	refreshDuration = 30 * time.Second
)

func main() {
	address := flag.String("address", "", "The server address")
	enableTLS := flag.Bool("tls", false, "Enable mutual TLS")

	flag.Parse()
	log.Printf("Dial server %s, TLS = %t", *address, *enableTLS)

	transportOption := grpc.WithInsecure() // no TLS by default

	if *enableTLS {
		tlsCredentials, err := loadTlsCred()
		if err != nil {
			log.Fatal("Cannot load TLS credentials: ", err)
		}
		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}

	// Connection for the auth interceptor.
	cc1, err := grpc.Dial(*address, transportOption)
	if err != nil {
		log.Fatal("Cannot dial server: ", err)
	}

	// Auth.
	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("Cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(
		*address,
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("Cannot dial server: ", err)
	}

	// Creating the RPC client.
	laptopClient := client.NewLaptopClient(cc2)
	// laptopClient := pb.NewLaptopServiceClient(conn)

	// testCreateLaptop(laptopClient)
	// testSearchLaptop(laptopClient)
	// testUploadImage(laptopClient)
	testRateLaptop(laptopClient)
}
