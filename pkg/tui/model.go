package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/avinash-apk/sentinel/pkg/bus" // replace with your module
)

// styles for the ui
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	eventStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))
)

type Model struct {
	events []string       // list of logs to show
	sub    chan bus.Event // channel to listen for new events
}

// initial state of the ui
func InitialModel(sub chan bus.Event) Model {
	return Model{
		events: []string{},
		sub:    sub,
	}
}

func (m Model) Init() tea.Cmd {
	// start the listener loop immediately
	return waitForActivity(m.sub)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	// handle key presses
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	// handle incoming events from the bus
	case bus.Event:
		// format the log line
		logLine := fmt.Sprintf("[%s] %v", msg.Topic, msg.Payload)
		m.events = append(m.events, logLine)
		
		// keep only last 10 events to save space
		if len(m.events) > 10 {
			m.events = m.events[1:]
		}
		
		// wait for the next event
		return m, waitForActivity(m.sub)
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}

	// render title
	s.WriteString(titleStyle.Render("SENTINEL DASHBOARD") + "\n\n")

	// render event log
	for _, e := range m.events {
		s.WriteString(eventStyle.Render(e) + "\n")
	}

	// render footer
	s.WriteString("\npress 'q' to quit\n")

	return s.String()
}

// this function converts a channel receive into a bubble tea message
func waitForActivity(sub chan bus.Event) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}