package ingest

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/avinash-apk/sentinel/pkg/bus" // Match your module name
)

type DiscordIngestor struct {
	Session *discordgo.Session
	Bus     *bus.EventBus
}

func NewDiscordIngestor(token string, b *bus.EventBus) (*DiscordIngestor, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &DiscordIngestor{Session: s, Bus: b}, nil
}

func (d *DiscordIngestor) Start() error {
	// Register the messageCreate func as a callback for MessageCreate events.
	d.Session.AddHandler(d.messageCreate)

	// Open a websocket connection to Discord and begin listening.
	d.Session.Identify.Intents = discordgo.IntentsGuildMessages
	err := d.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	fmt.Println("🎧 Discord Listener is Active")
	return nil
}

// This function is called every time a new message is created on any channel that the authenticated bot has access to.
func (d *DiscordIngestor) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Create a structured payload
	payload := map[string]string{
		"platform": "discord",
		"id":       m.ChannelID, // <--- THIS IS THE MAGIC ID YOU NEED
		"user":     m.Author.Username,
		"message":  m.Content,
	}

	// Publish to the bus
	d.Bus.Publish("discord:message", payload)
}