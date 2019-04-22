package main

import (
	"context"
	"log"

	"../gchat"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	defer conn.Close()

	c := gchat.NewChatServiceClient(conn)
	resp, err := c.Register(context.Background(), &gchat.UserContent{Username: "mange", Password: "passwd"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Response from server: %s", resp.Token)

}
