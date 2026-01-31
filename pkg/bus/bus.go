package bus

import (
	"sync"
)

// Event is the generic packet of data moving through our system
type Event struct {
	Topic   string
	Payload interface{} // interface{} means "any type of data"
}

// DataChannel is a pipe that carries Events
type DataChannel chan Event

// DataChannelSlice is just a list of those pipes
type DataChannelSlice []DataChannel

// EventBus is our Air Traffic Controller
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex // Mutex protects us from race conditions
}

// NewEventBus creates the controller
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string]DataChannelSlice),
	}
}

// Subscribe lets a listener say "I want to know about 'topic'"
func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	defer eb.rm.Unlock()
	// Append the new channel to the list of subscribers for this topic
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
}

// Publish broadcasts a message to everyone listening to 'topic'
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	defer eb.rm.RUnlock()

	// If no one is listening, do nothing
	if chans, found := eb.subscribers[topic]; found {
		// Create the event object
		go func(event Event, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				// Send the data into the pipe
				ch <- event
			}
		}(Event{Topic: topic, Payload: data}, chans)
	}
}