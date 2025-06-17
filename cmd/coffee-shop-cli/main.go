package main

import (
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
)

type app struct {
	response   string
	views      []tea.Model
	activeView int
}

func newApp() app {
	views := []tea.Model{outsideView{}, insideView{}}
	activeView := 0
	return app{views: views, activeView: activeView}
}

func (a app) Init() tea.Cmd {
	var cmd tea.Cmd
	for _, view := range a.views {
		cmd = tea.Batch(cmd, view.Init())
	}
	return cmd
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			if a.activeView == len(a.views)-1 {
				a.activeView = 0
			} else {
				a.activeView += 1
			}
			return a, nil
		}
	}
	view, cmd := a.views[a.activeView].Update(msg)
	a.views[a.activeView] = view
	return a, cmd
}

func (a app) View() string {
	return a.views[a.activeView].View()
}

func main() {
	f, err := tea.LogToFile("coffee-shop-cli.log", "coffee-shop-cli")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	prog := tea.NewProgram(newApp())
	_, err = prog.Run()
	if err != nil {
		panic(err)
	}
}
