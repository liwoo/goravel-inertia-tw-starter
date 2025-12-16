package events

import (
	"github.com/goravel/framework/contracts/event"
)

// ApplicationApproved is fired when an application is approved by an admin
type ApplicationApproved struct {
}

func (e *ApplicationApproved) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
