package events

import (
	"github.com/goravel/framework/contracts/event"
)

// ApplicationRejected is fired when an application is rejected by an admin
type ApplicationRejected struct {
}

func (e *ApplicationRejected) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
