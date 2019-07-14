package handler

import (
	"context"
	"fmt"
	"log"
	"io"
	"errors"

	"google.golang.org/grpc"
	"github.com/magnusbrattlof/go-grpc-chat/gchat"
)

var (
	client gchat.ChatServiceClient
	sequence int32 = 0
	token string
)

type UserData struct {
	Username string
	Token string
	Chat []string
}

// TODO implement defer conn.Close()
func init() {
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	client = gchat.NewChatServiceClient(conn)
}

func GetChats() {
	stream, err := client.GetChats(context.Background(), &gchat.RegisterResponse{Token: token})
	if err != nil {
		fmt.Printf("Error")
	}
	for {
		chat, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Epic faiiiil")
		}
		fmt.Println(chat.ChatName)
	}
}

func SendMessage(message string, token string) (error) {

	if token == "" {
		return errors.New("Token is empty, have you singed in?")
	}

	_, err := client.SendMessage(context.Background(), &gchat.ChatMessage{Sequence: sequence, Message: message})

	if (err != nil) {
		log.Fatalf("Error")
	}
	sequence++
	return nil
}

func Register() (*UserData, error) {
	var username, password string

	fmt.Printf("Username: ")
	fmt.Scanf("%s", &username)

	fmt.Printf("Password: ")
	fmt.Scanf("%s", &password)

	response, err := client.Register(context.Background(), &gchat.UserContent{Username: username, Password: password})
	token = response.Token

	if err != nil {
		log.Fatalf("Here is an error: %v", err)
		return nil, err
	} else {
		log.Printf("You have successfully registerd")
		return &UserData{Username: username, Token: token}, nil
	}
}
