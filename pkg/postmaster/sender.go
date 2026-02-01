package postmaster

// Messenger defines the behavior for anything that can send a message
type Messenger interface {
	Send(destination string, body string) error
}