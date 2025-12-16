package events

import (
	"github.com/goravel/framework/contracts/event"
)

// EventCreated is fired when a new event is created by an admin
type EventCreated struct {
}

func (e *EventCreated) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
