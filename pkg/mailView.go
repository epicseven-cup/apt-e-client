package pkg

import (
	"encoding/base64"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/gmail/v1"
)

type MailView struct {
	message *gmail.Message
	root    tea.Model
}

func (m MailView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m MailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return m.root, nil
		}
	}

	// Note that we're not returning a command.
	return m, nil
}

func (m MailView) View() string {
	//s := fmt.Sprintf("Header: %s\n", m.message.Payload.Headers)
	header := ""
	for i := range m.message.Payload.Headers {
		header += fmt.Sprintf("Name: %s\nValue:%s\n", m.message.Payload.Headers[i].Name, m.message.Payload.Headers[i].Value)
	}
	s := fmt.Sprintf("Header: %s\n", header)
	for _ = range s {
		s += "-"
	}

	content := make([]byte, 1)
	for i := range m.message.Payload.Parts {
		currentFrame, err := base64.URLEncoding.DecodeString(m.message.Payload.Parts[i].Body.Data)
		if err != nil {
			fmt.Sprintf("Error: %v", err)
		}
		content = append(content, currentFrame...)
	}

	s += fmt.Sprintf("\nBody:\n\n%s", string(content))

	// The footer
	s += "\nPress q to quit.\n"
	return s
}
