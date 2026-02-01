package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// REPLACE THESE with your actual module path
	"github.com/avinash-apk/sentinel/pkg/bus"
	"github.com/avinash-apk/sentinel/pkg/postmaster"
)

// STYLES
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

// STATE MANAGEMENT
type sessionState int

const (
	viewMode sessionState = iota
	replyMode
)

// DATA STRUCTURE FOR LIST ITEMS
type Notification struct {
	Platform string
	ID       string // The Channel ID or Issue Num
	User     string
	Message  string
}

type Model struct {
	state         sessionState
	notifications []Notification
	cursor        int             // Which item is selected
	textInput     textinput.Model // The reply box
	sub           chan bus.Event
	
	// Senders
	discordSender *postmaster.DiscordSender
	slackSender   *postmaster.SlackSender
}

func InitialModel(sub chan bus.Event, ds *postmaster.DiscordSender, ss *postmaster.SlackSender) Model {
	ti := textinput.New()
	ti.Placeholder = "Type your reply..."
	ti.CharLimit = 156
	ti.Width = 30

	return Model{
		state:         viewMode,
		notifications: []Notification{},
		sub:           sub,
		textInput:     ti,
		discordSender: ds,
		slackSender:   ss,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, waitForActivity(m.sub))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Handle Incoming Events
	case bus.Event:
		// Convert map payload to Notification struct
		if data, ok := msg.Payload.(map[string]string); ok {
			notif := Notification{
				Platform: data["platform"],
				ID:       data["id"],
				User:     data["user"],
				Message:  data["message"],
			}
			// Prepend (add to top)
			m.notifications = append([]Notification{notif}, m.notifications...)
		}
		return m, waitForActivity(m.sub)

	// Handle Keypresses
	case tea.KeyMsg:
		switch m.state {

		// VIEW MODE: Navigate the list
		case viewMode:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.notifications)-1 {
					m.cursor++
				}
			case "enter":
				// Switch to Reply Mode
				if len(m.notifications) > 0 {
					m.state = replyMode
					m.textInput.Focus()
				}
			}

		// REPLY MODE: Type and Send
		case replyMode:
			switch msg.String() {
			case "esc":
				m.state = viewMode
				m.textInput.Blur()
				m.textInput.Reset()
			case "enter":
				// EXECUTE THE REPLY
				target := m.notifications[m.cursor]
				replyText := m.textInput.Value()

				// Send via Postmaster based on platform
				if target.Platform == "discord" && m.discordSender != nil {
					m.discordSender.Send(target.ID, replyText)
				} else if target.Platform == "slack" && m.slackSender != nil {
					m.slackSender.Send(target.ID, replyText)
				}

				// Reset UI
				m.state = viewMode
				m.textInput.Reset()
				return m, nil
			}
		}
	}

	// Update Text Input bubble if in reply mode
	if m.state == replyMode {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString("SENTINEL COMMAND CENTER\n\n")

	// RENDER LIST
	for i, n := range m.notifications {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // selected
		}

		// Check if selected
		style := noStyle
		if m.cursor == i {
			style = focusedStyle
		}

		line := fmt.Sprintf("%s [%s] %s: %s", cursor, n.Platform, n.User, n.Message)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n")

	// RENDER REPLY BOX
	if m.state == replyMode {
		s.WriteString(fmt.Sprintf("Replying to %s:\n", m.notifications[m.cursor].User))
		s.WriteString(m.textInput.View())
		s.WriteString("\n(Enter to send, Esc to cancel)")
	} else {
		s.WriteString("(Press Enter on a message to reply)")
	}

	return s.String()
}

func waitForActivity(sub chan bus.Event) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}