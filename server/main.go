package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net"

	"../gchat"
	"google.golang.org/grpc"
)

type Server struct {
}

func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := Server{}
	grpcServer := grpc.NewServer()
	gchat.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}

func (s *Server) Register(ctx context.Context, in *gchat.UserContent) (*gchat.RegisterResponse, error) {
	fmt.Printf("Your username was: %s and password: %s", in.Username, in.Password)
	return &gchat.RegisterResponse{Token: tokenGenerator()}, nil
}

func (s *Server) Login(ctx context.Context, in *gchat.UserContent) (*gchat.LoginResponse, error) {
	fmt.Printf("Your username was: %s and password: %s", in.Username, in.Password)
	return &gchat.LoginResponse{Token: "hello"}, nil
}

func (s *Server) Logout(ctx context.Context, in *gchat.UserContent) (*gchat.LogoutResponse, error) {
	fmt.Printf("Your username was: %s and password: %s", in.Username, in.Password)
	return &gchat.LogoutResponse{Token: "hello"}, nil
}

// Returns a string with random chars
func tokenGenerator() string {
	// Create the slice with
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
