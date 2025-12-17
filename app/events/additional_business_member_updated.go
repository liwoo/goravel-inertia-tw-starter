package events

import (
	"github.com/goravel/framework/contracts/event"
)

// AdditionalBusinessMemberUpdated event is fired when an additional business member is updated
type AdditionalBusinessMemberUpdated struct {
}

func (e *AdditionalBusinessMemberUpdated) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
