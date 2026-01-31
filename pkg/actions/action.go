package actions

// action is anything that can be executed
type Action interface {
	Execute(payload interface{}) error
}