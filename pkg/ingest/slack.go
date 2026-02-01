package ingest

import (
	"fmt"
	"sync" // <--- Added for thread safety

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	// REPLACE THIS with your actual module path
	"github.com/avinash-apk/sentinel/pkg/bus"
)

type SlackIngestor struct {
	Api          *slack.Client
	SocketClient *socketmode.Client
	Bus          *bus.EventBus
	UserCache    map[string]string
	CacheMu      sync.Mutex // <--- Lock to prevent crashes
	MyBotUserID  string
}

func NewSlackIngestor(appToken string, botToken string, b *bus.EventBus) *SlackIngestor {
	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)

	// Fetch Bot ID for self-ignore
	authTest, err := api.AuthTest()
	myID := ""
	if err != nil {
		fmt.Printf("⚠️ Warning: Could not fetch Bot ID: %v\n", err)
	} else {
		myID = authTest.UserID
	}

	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
	)

	return &SlackIngestor{
		Api:          api,
		SocketClient: client,
		Bus:          b,
		UserCache:    make(map[string]string),
		MyBotUserID:  myID,
	}
}

// Thread-safe name lookup
func (s *SlackIngestor) getUserName(userID string) string {
	s.CacheMu.Lock()
	val, ok := s.UserCache[userID]
	s.CacheMu.Unlock()

	if ok {
		return val
	}

	// Fetch from API (Slow Network Call)
	user, err := s.Api.GetUserInfo(userID)
	if err != nil {
		return userID
	}

	s.CacheMu.Lock()
	s.UserCache[userID] = user.RealName
	s.CacheMu.Unlock()
	
	return user.RealName
}

func (s *SlackIngestor) Start() {
	fmt.Println("Slack Listener is Active")

	go func() {
		for evt := range s.SocketClient.Events {
			switch evt.Type {

			case socketmode.EventTypeConnected:
				// Connected

			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}

				// 1. ACKNOWLEDGE IMMEDIATELY
				// This tells Slack "I got it" so they don't resend it.
				s.SocketClient.Ack(*evt.Request)

				// 2. PROCESS ASYNC
				// We launch a background worker so the main loop isn't blocked.
				go func(e slackevents.EventsAPIEvent) {
					switch e.Type {
					case slackevents.CallbackEvent:
						innerEvent := e.InnerEvent
						switch ev := innerEvent.Data.(type) {
						case *slackevents.MessageEvent:
							
							// Filter out Bot messages
							if ev.BotID != "" {
								return
							}
							if ev.User == s.MyBotUserID {
								return
							}

							// This API call takes time, but now it won't block the next Ack
							realName := s.getUserName(ev.User)

							payload := map[string]string{
								"platform": "slack",
								"id":       ev.Channel,
								"user":     realName,
								"message":  ev.Text,
							}
							s.Bus.Publish("slack:message", payload)
						}
					}
				}(eventsAPIEvent) // Pass event into the closure
			}
		}
	}()

	s.SocketClient.Run()
}