package postmaster

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DiscordSender struct {
	Session *discordgo.Session
}

// NewDiscordSender authenticates with the bot token
func NewDiscordSender(token string) (*DiscordSender, error) {
	// create a session (does not connect yet)
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	
	return &DiscordSender{Session: sess}, nil
}

// Send posts a message to a specific channel id
func (d *DiscordSender) Send(destination string, body string) error {
	// open the connection (required for some state, though simple sends can sometimes work via rest)
	// for a cli, we open, send, and defer close
	if err := d.Session.Open(); err != nil {
		return fmt.Errorf("error opening discord session: %v", err)
	}
	defer d.Session.Close()

	// destination is the channel id (e.g., "123456789")
	_, err := d.Session.ChannelMessageSend(destination, body)
	if err != nil {
		return fmt.Errorf("discord send failed: %v", err)
	}

	return nil
}