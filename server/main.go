package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/magnusbrattlof/go-grpc-chat/gchat"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

/* Defining all types used in the server */
type User struct {
	username *string
	password *string
	token    string
}

type ChatMessage struct {
	username  *string
	sequence  int32
	timestamp time.Time
	message   string
}

type ChatCache struct {
	Name     string
	ID       int32
	Messages []*ChatMessage
	Clients  map[string]string
}

type Server struct {
}

/* Defining all global variables and lists */
var (
	database []User
	chats    []*gchat.Chats
	ChatList []*ChatCache
)

func init() {
	// Initialize som dummy chats
	chats = append(chats, &gchat.Chats{ChatName: "family", ChatID: 0})
	chats = append(chats, &gchat.Chats{ChatName: "work", ChatID: 1})
	chats = append(chats, &gchat.Chats{ChatName: "personal", ChatID: 2})
	for _, chat := range chats {
		ChatList = append(ChatList, &ChatCache{Name: chat.ChatName, ID: chat.ChatID, Clients: make(map[string]string)})
	}
}

func (s *Server) Register(ctx context.Context, in *gchat.UserContent) (*gchat.RegisterResponse, error) {

	if userNameExists(in.Username) {
		return nil, errors.New("Username already exists...")
	}

	token := tokenGenerator()
	database = append(database, User{username: &in.Username, password: &in.Password, token: token})
	return &gchat.RegisterResponse{Token: token}, nil
}

func (s *Server) GetChats(in *gchat.RegisterResponse, stream gchat.ChatService_GetChatsServer) error {
	for _, chat := range chats {
		if err := stream.Send(chat); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) CreateChat(ctx context.Context, in *gchat.Chats) (*gchat.Response, error) {
	//
	return nil, nil
}

func (s *Server) Login(ctx context.Context, in *gchat.UserContent) (*gchat.LoginResponse, error) {

	ChatList[2].Clients[in.Username] = fmt.Sprintf(":%d", in.ReceiverPort)
	fmt.Println("Added port ", in.ReceiverPort)
	// if userNameExists(in.Username) {
	// 	return nil, errors.New("Username already exists...")
	// }

	fmt.Printf("Your username was: %s and password: %s\n", in.Username, in.Password)
	return &gchat.LoginResponse{Token: "hello"}, nil
}

func (s *Server) Logout(ctx context.Context, in *gchat.UserContent) (*gchat.LogoutResponse, error) {
	fmt.Printf("Your username was: %s and password: %s\n", in.Username, in.Password)
	return &gchat.LogoutResponse{Token: "hello"}, nil
}

func (s *Server) SendMessage(ctx context.Context, in *gchat.ChatMessage) (*gchat.MessageResponse, error) {

	fmt.Println("Received message: ")
	fmt.Println("--------------------")
	fmt.Printf("Msg:\t%s\nSeq:\t%d\nChat:\t%s\nTime:\t%s\nUname:\t%s\n", in.Msg, in.Sequence, in.Chat, in.Timestamp, in.Username)
	// Returns the list of all the saved messages
	deliverMessage(in, in.Username)
	return &gchat.MessageResponse{Val: true}, nil
}

func deliverMessage(in *gchat.ChatMessage, sender string) {

	for uname, port := range ChatList[2].Clients {
		fmt.Println("Port ", port)
		if sender != uname {
			fmt.Println("Sending...")
			conn, err := grpc.Dial(port, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Error connecting: %v", err)
			}
			client := gchat.NewReceiverClient(conn)
			client.ReceiveMessage(context.Background(), in)
			fmt.Println("Just sent message...")
			defer conn.Close()
		}
	}
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

func userNameExists(username string) bool {
	for _, i := range database {
		if username == *i.username {
			return true
		}
	}
	return false
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
