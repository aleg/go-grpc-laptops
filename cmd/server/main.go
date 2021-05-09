package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/aleg/go-grpc-laptops/service"
	"github.com/aleg/go-grpc-laptops/stores"
	"github.com/aleg/go-grpc-laptops/users"
	"google.golang.org/grpc"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func main() {
	port := flag.Int("port", 0, "The server port")
	enableTLS := flag.Bool("tls", false, "Enable mutual TLS")

	flag.Parse()
	log.Printf("Start server on port %d, TLS = %t", *port, *enableTLS)

	// Creating some users and the auth server.
	userStore := stores.NewInMemoryUserStore()
	createUsers(userStore)
	jwtManager := users.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := stores.NewInMemoryLaptopStore()
	imageStore := stores.NewDiskImageStore("tmp/uploaded-img")
	ratingStore := stores.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	// Interceptors.
	authInterceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	unaryInterceptor := grpc.UnaryInterceptor(authInterceptor.Unary())
	streamInterceptor := grpc.StreamInterceptor(authInterceptor.Stream())
	serverOpts := []grpc.ServerOption{unaryInterceptor, streamInterceptor}

	// TLS
	if *enableTLS {
		tlsCred, err := loadTlsCred()
		if err != nil {
			log.Fatal("Cannot load TLS credentials: ", err)
		}
		serverOpts = append(serverOpts, grpc.Creds(tlsCred))
	}

	// Creating the gRPC server with the given options and
	// registering the services.
	grpcServer := grpc.NewServer(serverOpts...)
	// unaryInterceptor := grpc.UnaryInterceptor(unaryInterceptorHandler)
	// streamInterceptor := grpc.StreamInterceptor(streamInterceptorHandler)
	// grpcServer := grpc.NewServer(unaryInterceptor, streamInterceptor)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	// Serving.
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot serve from the server: ", err)
	}
}

// func unaryInterceptorHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//         log.Println("--> Unary interceptor: ", info.FullMethod)
//         return handler(ctx, req)
// }

// func streamInterceptorHandler(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
//         log.Println("--> Stream interceptor: ", info.FullMethod)
//         return handler(srv, stream)
// }
