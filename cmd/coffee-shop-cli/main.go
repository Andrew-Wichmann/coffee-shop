package main

import (
	"io"
	"log"
	"net/http"

	_ "github.com/charmbracelet/bubbles"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
)

type app struct {
	response string
}

type requestError error
type requestResp string

func (a app) Init() tea.Cmd {
	return nil
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
	log.Println("sending http request")
	resp, err := http.Get("http://localhost:8080/healthcheck")
	if err != nil {
		return requestError(err)
	}
	defer resp.Body.Close()
	buf := make([]byte, 1000)
	_, err = resp.Body.Read(buf)
	if err != io.EOF {
		return requestError(err)
	}
	return requestResp(buf)
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
