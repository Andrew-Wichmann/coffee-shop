package main

import (
	"log"

	_ "github.com/charmbracelet/bubbles"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type app struct {
	response string
}

var conn *websocket.Conn

type requestError error
type requestResp string

func createConnection() tea.Msg {
	log.Println("Connecting to websocket")
	var err error
	conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		return requestError(err)
	}
	return nil
}

func (a app) Init() tea.Cmd {
	return createConnection
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			return a, sendRequest
		}
	}
	if response, ok := msg.(requestResp); ok {
		a.response = string(response)
	}
	if response, ok := msg.(requestError); ok {
		a.response = string(response.Error())
	}
	return a, nil
}

func (a app) View() string {
	if a.response != "" {
		return a.response
	}
	return "Hello world"
}

func sendRequest() tea.Msg {
	err := conn.WriteMessage(websocket.TextMessage, []byte("enter"))
	if err != nil {
		return requestError(err)
	}
	_, _msg, err := conn.ReadMessage()
	if err != nil {
		return requestError(err)
	}
	log.Println(string(_msg))
	return requestResp(_msg)
}

func main() {
	f, err := tea.LogToFile("coffee-shop-cli.log", "coffee-shop-cli")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	prog := tea.NewProgram(app{})
	_, err = prog.Run()
	if err != nil {
		panic(err)
	}
}
