package engine

import (
	"github.com/avinash-apk/sentinel/pkg/actions" // replace with your module path
	"github.com/avinash-apk/sentinel/pkg/bus"     // replace with your module path
)

// rule connects a topic to an action
type Rule struct {
	Topic  string
	Action actions.Action
}

type Engine struct {
	Bus   *bus.EventBus
	Rules []Rule
}

func (e *Engine) Start() {
	// create a master channel to listen to everything
	// in a real app, we might filter, but for now we listen to common topics
	ch := make(chan bus.Event)
	
	// subscribe the engine to github events
	e.Bus.Subscribe("github:event", ch)

	// process events loop
	go func() {
		for event := range ch {
			for _, rule := range e.Rules {
				// simple match: does the topic match?
				if rule.Topic == event.Topic {
					// run the action in a goroutine so we don't block
					go rule.Action.Execute(event.Payload)
				}
			}
		}
	}()
}