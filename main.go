package main

import (
	"context"
	"fmt"
	"github.com/magnusbrattlof/go-grpc-chat/gchat"
	"github.com/magnusbrattlof/go-grpc-chat/gchat/handler"
	"github.com/marcusolsson/tui-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type post struct {
	username string
	message  string
	time     string
}

type Server struct {
}

var (
	err        error
	posts      = []post{}
	option     string
	userData   *handler.UserData
	history    *tui.Box
	ui         tui.UI
	msgChannel = make(chan *gchat.ChatMessage)
)

func main() {

	fmt.Println("==========> Welcome to go-grpc-chat <==========")
	userData, err = handler.Register()
	if err != nil {
		log.Println(err)
	}

	chats := handler.GetChats()
	fmt.Println("Available chats for you: ")
	for i, chat := range *chats {
		fmt.Printf("%d) %s\n", i, chat)
	}
	fmt.Println("Select a chat:")
	var opt int
	fmt.Scanf("%d", &opt)
	currentChat := (*chats)[opt]
	go runReceiver()
	gui_chat_handler(currentChat)
}

func runReceiver() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", userData.ReceiverPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := Server{}
	grpcServer := grpc.NewServer()
	gchat.RegisterReceiverServer(grpcServer, &s)

	log.Println("Server is running on port ", userData.ReceiverPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}

func gui_chat_handler(currentChat string) {

	history = tui.NewVBox()

	for _, m := range posts {
		history.Append(tui.NewHBox(
			tui.NewLabel(m.time),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.username))),
			tui.NewLabel(m.message),
			tui.NewSpacer(),
		))
	}

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		handler.SendMessage(e.Text(), userData.Token, currentChat, userData.Username)

		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", userData.Username))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})

	root := tui.NewHBox(chat)

	ui, err = tui.New(root)

	go func() {
		for msg := range msgChannel {
			// we need to make the change via ui update to make sure the ui is repaint correctly
			ui.Update(func() {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format("15:04")),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", msg.Username))),
					tui.NewLabel(msg.Msg),
					tui.NewSpacer(),
				))
				input.SetText("")
			})
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ReceiveMessage(ctx context.Context, in *gchat.ChatMessage) (*gchat.MessageResponse, error) {
	msgChannel <- in
	return &gchat.MessageResponse{Val: true}, nil
}
