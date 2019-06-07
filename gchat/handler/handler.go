package handler

import (
	"context"
	"fmt"
	"log"
	"errors"

	"github.com/magnusbrattlof/go-grpc-chat/gchat"
)

func Message_handler(c gchat.ChatServiceClient, token *string) (error) {
	var message string
	var sequence int32 = 0

	if token == nil {
		return errors.New("Token is empty, have you singed in?")
	}

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

func Register_handler(c gchat.ChatServiceClient) (*string, error) {
	var username, password string

	fmt.Printf("Username: ")
	fmt.Scanf("%s", &username)
	
	fmt.Printf("Password: ")
	fmt.Scanf("%s", &password)

	token, err := c.Register(context.Background(), &gchat.UserContent{Username: username, Password: password})

	if err != nil {
		log.Fatalf("Here is an error: %v", err)
		return nil, err
	} else {
		log.Printf("You have successfully registerd")
		return &token.Token, nil
	}
}