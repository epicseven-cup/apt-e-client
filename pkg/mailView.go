package pkg

import (
	"encoding/base64"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/gmail/v1"
	"strings"
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

		//case tea.WindowSizeMsg:
		//	m.height = msg.Height
		//	m.width = msg.Width
	}

	// Note that we're not returning a command.
	return m, nil
}

func (m MailView) View() string {
	header := map[string]string{}
	for i := range m.message.Payload.Headers {
		header[m.message.Payload.Headers[i].Name] = m.message.Payload.Headers[i].Value
	}
	s := fmt.Sprintf("Subject: %s\nFrom: %s\nTo:%s\n", header["Subject"], header["From"], header["To"])

	for _ = range 40 {
		s += "-"
	}

	s += fmt.Sprintf("\nDate: %s\n", header["Date"])

	body := map[string][]byte{}
	for i := range m.message.Payload.Parts {
		content, err := base64.URLEncoding.DecodeString(m.message.Payload.Parts[i].Body.Data)
		if err != nil {
			fmt.Errorf("Error: %v", err)
		}
		body[m.message.Payload.Parts[i].MimeType] = content
	}

	content := string(body["text/plain"])
	content = strings.TrimSpace(content)
	s += fmt.Sprintf("Body:\n\n%s\n", content)
	// The footer
	s += "Press [ESC] to go back.\n"
	s += "Press q to quit.\n"
	return s
}
