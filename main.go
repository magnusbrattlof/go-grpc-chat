package main


import (
	"time"
	"fmt"
	"log"
	"github.com/marcusolsson/tui-go"
	"github.com/magnusbrattlof/go-grpc-chat/gchat/handler"
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

var (
	option int
	err error
	userData *handler.UserData
)

func main() {

	fmt.Println("Welcome to go-grpc-chat\nWhat would you like to do?")

	for {

		fmt.Println("1) Register\n2) Chat\n3) Exit")
		fmt.Scanf("%d", &option)

		switch option {

		case 1:
			userData, err = handler.Register()

			if err != nil {
				fmt.Println("Error, could not register")
			}

		case 2:
			handler.GetChats()
			gui_chat_handler()
		}
	}
}

func gui_chat_handler() {

	if userData == nil {
		log.Fatalf("You have not signed in yet")
	}
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
		handler.SendMessage(e.Text(), userData.Token)

		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", userData.Username))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
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
