package postmaster

import (
	"github.com/slack-go/slack"
)

type SlackSender struct {
	Client *slack.Client
}

func NewSlackSender(token string) *SlackSender {
	return &SlackSender{
		Client: slack.New(token),
	}
}

func (s *SlackSender) Send(destination string, body string) error {
	// destination is the Channel ID (e.g., C12345)
	_, _, err := s.Client.PostMessage(destination, slack.MsgOptionText(body, false))
	return err
}