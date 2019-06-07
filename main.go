package main


import (
	"time"
	"fmt"
	"log"
	"github.com/marcusolsson/tui-go"
	"github.com/magnusbrattlof/go-grpc-chat/gchat"
	"github.com/magnusbrattlof/go-grpc-chat/gchat/handler"
	"google.golang.org/grpc"
)

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "john", message: "hi, what's up?", time: "14:41"},
	{username: "jane", message: "not much", time: "14:43"},
}

func main() {

	var (
		conn               *grpc.ClientConn
		option int
		token *string
		err error
	)

	token = nil
	conn, err = grpc.Dial(":7777", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	defer conn.Close()

	c := gchat.NewChatServiceClient(conn)

	
	fmt.Println("Welcome to go-grpc-chat\nWhat would you like to do?")

	for {

		fmt.Println("1) Register\n2) Chat\n3) Exit")
		fmt.Scanf("%d", &option)

		switch option {

		case 1:
			token, err = handler.Register_handler(c)
			
			if err != nil {
				fmt.Println("Error, could not register")
			}
		
		case 2:
			err = handler.Message_handler(c, token)

			if err != nil {
				fmt.Println(err)
			}

		}
	}
	
	//

}

func gui_chat_handler() {

	history := tui.NewVBox()

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
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
		// Here we can add functions that we want to execute whenever we have user input
	})

	root := tui.NewHBox(chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}