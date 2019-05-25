package main

import (
	"context"
	"fmt"
	"log"

	"../gchat"
	"google.golang.org/grpc"
)

func main() {
	var (
		conn               *grpc.ClientConn
		option             int
	)

	fmt.Println("Welcome to go-grpc-chat\nWhat would you like to do?")

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	defer conn.Close()

	c := gchat.NewChatServiceClient(conn)

	for {
		fmt.Println("1) \t Register")
		fmt.Println("2) \t Chat")
		fmt.Println("3) \t Exit")
		fmt.Scanf("%d", &option)

		switch option {
		
		case 1:
			token, err := register_handler(c)
		
		case 2:

		}
	}
}

func message_handler(c gchat.ChatServiceClient, token string) {
	var message string
	var sequence int32 = 0

	for {
		fmt.Printf("@chat$: ")
		fmt.Scanf("%d", &message)
		fmt.Println("@chat$: ")
		_, err := c.SendMessage(context.Background(), &gchat.ChatMessage{Sequence: sequence, Message: message})
	
		if (err != nil) {
			log.Fatalf("Error")
		}
	}
}

func register_handler(c gchat.ChatServiceClient) (*string, error) {
	var username, password string

	fmt.Printf("Username: ")
	fmt.Scanf("%s", &username)
	
	fmt.Printf("Password: ")
	fmt.Scanf("%s", &password)

	token, err := c.Register(context.Background(), &gchat.UserContent{Username: username, Password: password})

	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil, err
	} else {
		log.Printf("You have successfully registerd")
		return &token.Token, nil
	}
}