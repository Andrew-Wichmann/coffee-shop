package main

import (
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
)

type app struct {
	response string
	views    []tea.Model
}

func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			return a, nil
		}
	}
	return a, nil
}

func (a app) View() string {
	return "Hello world"
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
