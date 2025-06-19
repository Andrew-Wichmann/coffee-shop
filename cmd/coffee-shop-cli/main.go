package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("coffee-shop-cli.log", "coffee-shop-cli")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	prog := tea.NewProgram(initializeApp(), tea.WithAltScreen())
	_, err = prog.Run()
	if err != nil {
		panic(err)
	}
}
