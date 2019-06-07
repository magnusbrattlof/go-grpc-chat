package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/magnusbrattlof/go-grpc-chat/gchat"
	"google.golang.org/grpc"
)

type User struct {
	username *string
	password *string
	token    string
}

var database []User

type Server struct {
}

func (s *Server) Register(ctx context.Context, in *gchat.UserContent) (*gchat.RegisterResponse, error) {

	for _, i := range database {
		if in.Username == *i.username {
			return nil, errors.New("Username already exists...")
		}
	}

	token := tokenGenerator()

	database = append(database, User{username: &in.Username, password: &in.Password, token: token})
	
	fmt.Println(len(database))
	printUsers()
	return &gchat.RegisterResponse{Token: token}, nil
}

func (s *Server) Login(ctx context.Context, in *gchat.UserContent) (*gchat.LoginResponse, error) {
	fmt.Printf("Your username was: %s and password: %s", in.Username, in.Password)
	return &gchat.LoginResponse{Token: "hello"}, nil
}

func (s *Server) Logout(ctx context.Context, in *gchat.UserContent) (*gchat.LogoutResponse, error) {
	fmt.Printf("Your username was: %s and password: %s", in.Username, in.Password)
	return &gchat.LogoutResponse{Token: "hello"}, nil
}

func (s *Server) SendMessage(ctx context.Context, in *gchat.ChatMessage) (*gchat.MessageResponse, error) {

	fmt.Printf("Received message: %s in sequence: %d", in.Message, in.Sequence)

	return &gchat.MessageResponse{Val: true}, nil
}

func deliver() {

}

// Returns a string with random chars
func tokenGenerator() string {
	// Create the slice with
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func printUsers() {
	for _, i := range database {
		fmt.Println(*i.username)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := Server{}
	grpcServer := grpc.NewServer()
	gchat.RegisterChatServiceServer(grpcServer, &s)

	log.Println("Server is running...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
