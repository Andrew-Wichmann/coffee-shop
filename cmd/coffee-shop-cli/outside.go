package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type outsideView struct{}

func (o outsideView) Init() tea.Cmd {
	return nil
}

func (o outsideView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o outsideView) View() string {
	return "The view from the outside"
}
