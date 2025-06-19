package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type size int

const (
	TOO_SMALL = iota
	MEDIUM
)

type app struct {
	size            size
	userInput       textarea.Model
	chatArea        viewport.Model
	messages        []string
	borderStyle     lipgloss.Style
	viewportWidth   int
	viewportHeight  int
	containerWidth  int
	containerHeight int
}

var conn *websocket.Conn

type requestError error
type requestResp string

func createConnection() tea.Msg {
	log.Println("Connecting to websocket")
	var err error
	// TODO: put api host in config
	conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		return requestError(err)
	}
	return nil
}

func initializeApp() app {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.Focus()
	ta.MaxHeight = 1
	vp := viewport.New(0, 0)
	vp.SetContent("")
	style := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
	return app{userInput: ta, chatArea: vp, borderStyle: style}
}

func (a app) Init() tea.Cmd {
	return tea.Batch(createConnection)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			request := a.userInput.Value()
			a.userInput.SetValue("")
			return a, func() tea.Msg { return sendRequest(request) }
		}
	}
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		switch {
		case msg.Width < 150 && msg.Height < 100:
			a.size = TOO_SMALL
			return a, nil
		default:
			a.size = MEDIUM
			a.containerWidth = 125
			a.containerHeight = 40
			a.chatArea.Height = a.containerHeight
			a.chatArea.Width = a.containerWidth
			a.userInput.SetWidth(a.containerWidth)
			a.userInput.SetHeight(1)
		}
		a.viewportWidth = msg.Width
		a.viewportHeight = msg.Height
	}
	if response, ok := msg.(requestResp); ok {
		a.messages = append(a.messages, string(response))
		a.chatArea.SetContent(strings.Join(a.messages, "\n"))
		a.chatArea.GotoBottom()
	}
	if response, ok := msg.(requestError); ok {
		log.Printf("Error received from the server %s", response.Error())
	}
	var cmd tea.Cmd
	a.userInput, cmd = a.userInput.Update(msg)
	return a, cmd
}

func (a app) View() string {
	if a.size == TOO_SMALL {
		return "Your terminal is too small to render this app. Please increase the terminal size."
	}
	return lipgloss.Place(a.viewportWidth, a.viewportHeight, lipgloss.Center, lipgloss.Center, a.borderStyle.Render(lipgloss.JoinVertical(lipgloss.Left, a.chatArea.View(), a.userInput.View())))
}

func sendRequest(request string) tea.Msg {
	err := conn.WriteMessage(websocket.TextMessage, []byte(request))
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
