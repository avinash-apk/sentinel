package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ID       string
	User     string
	Message  string
}

type Model struct {
	state         sessionState
	notifications []Notification
	cursor        int
	textInput     textinput.Model
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
		if data, ok := msg.Payload.(map[string]string); ok {
			notif := Notification{
				Platform: data["platform"],
				ID:       data["id"],
				User:     data["user"],
				Message:  data["message"],
			}
			m.notifications = append([]Notification{notif}, m.notifications...)
		}
		return m, waitForActivity(m.sub)

	// Handle Keypresses
	case tea.KeyMsg:
		switch m.state {

		// VIEW MODE
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
				if len(m.notifications) > 0 {
					m.state = replyMode
					m.textInput.Focus()
				}
			}

		// REPLY MODE
		case replyMode:
			switch msg.String() {
			case "esc":
				m.state = viewMode
				m.textInput.Blur()
				m.textInput.Reset()
			case "enter":
				target := m.notifications[m.cursor]
				replyText := m.textInput.Value()

				if target.Platform == "discord" && m.discordSender != nil {
					m.discordSender.Send(target.ID, replyText)
				} else if target.Platform == "slack" && m.slackSender != nil {
					m.slackSender.Send(target.ID, replyText)
				}

				m.state = viewMode
				m.textInput.Reset()
				return m, nil
			}
		}
	}

	if m.state == replyMode {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString("SENTINEL COMMAND CENTER\n\n")

	for i, n := range m.notifications {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		style := noStyle
		if m.cursor == i {
			style = focusedStyle
		}

		line := fmt.Sprintf("%s [%s] %s: %s", cursor, n.Platform, n.User, n.Message)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n")

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