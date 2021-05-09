package service

import (
	"context"
	"fmt"
	"log"

	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/aleg/go-grpc-laptops/stores"
	"github.com/aleg/go-grpc-laptops/users"
	"google.golang.org/grpc/codes"
)

type AuthServer struct {
	userStore  stores.UserStore
	jwtManager *users.JWTManager
}

func NewAuthServer(userStore stores.UserStore, jwtManager *users.JWTManager) *AuthServer {
	return &AuthServer{userStore: userStore, jwtManager: jwtManager}
}

// Login logs a user in and returns an auth token.
func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Received login-request for user %s", req.GetUsername())

	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, logError(err, codes.Internal, "Cannot find user")
	}
	if user == nil {
		msg := fmt.Sprintf("User \"%s\" not found", req.GetUsername())
		return nil, logError(err, codes.NotFound, msg)
	}
	if !user.IsCorrectPassword(req.GetPassword()) {
		return nil, logError(err, codes.NotFound, "Password not correct")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, logError(err, codes.Internal, "Cannot generate auth token")
	}

	log.Printf("Successfully generated the auth token for user %s", req.GetUsername())
	res := &pb.LoginResponse{
		AccessToken: token,
	}
	return res, nil
}
