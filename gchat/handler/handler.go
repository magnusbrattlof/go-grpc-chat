package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	// "os"
	"time"

	"github.com/magnusbrattlof/go-grpc-chat/gchat"
	"google.golang.org/grpc"
)

var (
	client   gchat.ChatServiceClient
	sequence int32 = 0
	token    string
	option   string
	USERNAME string = "GCHAT_USERNAME"
	PASSWORD string = "GCHAT_PASSWORD"
)

type UserData struct {
	Username     string
	Token        string
	Chat         []string
	ReceiverPort int32
}

// TODO implement defer conn.Close()
func init() {

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	client = gchat.NewChatServiceClient(conn)
}

// ADD BUFFERED CHATS DO MINIMIZE SERVER REQUESTS
func GetChats() *[]string {
	var chats []string
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
		chats = append(chats, chat.ChatName)
	}
	return &chats
}

func SendMessage(message, token, chat, username string) error {

	if token == "" {
		return errors.New("Token is empty, have you singed in?")
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	_, err := client.SendMessage(context.Background(), &gchat.ChatMessage{
		Sequence:  sequence,
		Msg:       message,
		Chat:      chat,
		Timestamp: t,
		Username:  username,
	})

	if err != nil {
		log.Fatalf("Error")
	}
	sequence++
	return nil
}

// func Login() (*UserData, error) {
// 	//username := os.Getenv(USERNAME)
// 	//password := os.Getenv(PASSWORD)

// 	// if username == "" || password == "" {
// 	// 	fmt.Println("No user settings were found in gchat environment variables...")
// 	// 	fmt.Println("Do you want to create a new account?")
// 	// 	fmt.Printf("Y/n ")
// 	// 	fmt.Scanln("%s", &option)
// 	// 	if option == "n" {
// 	// 		return nil, nil
// 	// 	} else if option == "Y" || option == "" {
// 	// 		Register()
// 	// 	}
// 	// }
// 	Register()
// 	rand.Seed(time.Now().UTC().UnixNano())
// 	rp := rand.Intn(100) + 6000
// 	receiverPort := int32(rp)
// 	response, err := client.Login(context.Background(), &gchat.UserContent{Username: username, Password: password, ReceiverPort: receiverPort})
// 	if err != nil {
// 		log.Printf("Could not login... %v\n", err)
// 		return nil, err
// 	}

// 	token = response.Token
// 	return &UserData{Username: username, Token: token, ReceiverPort: receiverPort}, nil
// }

func Register() (*UserData, error) {
	var username, password string

	fmt.Printf("Username: ")
	fmt.Scanf("%s", &username)

	fmt.Printf("Password: ")
	fmt.Scanf("%s", &password)

	response, err := client.Register(context.Background(), &gchat.UserContent{Username: username, Password: password})

	if err != nil {
		log.Fatalf("Here is an error: %v", err)
		return nil, err
	} else {
		token = response.Token
		log.Printf("You have successfully registerd")

		rand.Seed(time.Now().UTC().UnixNano())
		rp := rand.Intn(100) + 6000
		receiverPort := int32(rp)

		_, err := client.Login(context.Background(), &gchat.UserContent{Username: username, Password: password, ReceiverPort: receiverPort})
		if err != nil {
			log.Printf("Could not login... %v\n", err)
			return nil, err
		}
		return &UserData{Username: username, Token: token, ReceiverPort: receiverPort}, nil
	}
}
