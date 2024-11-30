package pkg

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	inbox  *Inbox
	cursor int // which to-do list item our cursor is pointing at
}

const inboxSize = 10

func InitialModel() model {
	inbox := InitInbox(inboxSize)
	//inbox := InitFakeInbox(inboxSize)
	return model{
		inbox:  inbox,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.inbox.messages)-1 {
				m.cursor++
			}

		// The "enter" key and the spacer (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter":
			return MailView{
				message: m.inbox.messages[m.cursor],
				root:    m,
			}, nil

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Here is your emails\n\n"
	// Iterate over our choices
	for index, message := range m.inbox.messages {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == index {
			cursor = ">" // cursor!
		}
		// Render the row
		s += fmt.Sprintf("%s | %s\n", cursor, message.Snippet)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
