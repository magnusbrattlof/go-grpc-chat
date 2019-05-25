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
		username, password string
		option             int
		token              string
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
			fmt.Printf("Username: ")
			fmt.Scanf("%s", &username)
			fmt.Printf("Password: ")
			fmt.Scanf("%s", &password)

			token, err := c.Register(context.Background(), &gchat.UserContent{Username: username, Password: password})
			if err != nil {
				log.Fatalf("Error: %v", err)
			} else {
				log.Printf("You have successfully registerd")
			}
		case 2:

		}
	}
}
