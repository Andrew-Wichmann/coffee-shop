package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type insideView struct{}

func (o insideView) Init() tea.Cmd {
	return nil
}

func (o insideView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o insideView) View() string {
	return "The view from the inside"
}
