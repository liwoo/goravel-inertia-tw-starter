package events

import (
	"github.com/goravel/framework/contracts/event"
)

// AdditionalBusinessMemberCreated event is fired when a new additional business member is created
type AdditionalBusinessMemberCreated struct {
}

func (e *AdditionalBusinessMemberCreated) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
