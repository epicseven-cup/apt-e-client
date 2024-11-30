package main

import (
	"apt-e-client/pkg"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	//f, err := tea.LogToFile("debug.log", "debug")
	//if err != nil {
	//	fmt.Println("fatal:", err)
	//	os.Exit(1)
	//}
	//defer f.Close()
	p := tea.NewProgram(pkg.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, There is an error: %v\n", err)
		os.Exit(1)
	}
}
