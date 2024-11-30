package main

import (
	"apt-e-client/pkg"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	p := tea.NewProgram(pkg.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, There is an error: %v", err)
		os.Exit(1)
	}
}
