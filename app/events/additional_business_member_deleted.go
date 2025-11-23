package events

import (
	"github.com/goravel/framework/contracts/event"
)

// AdditionalBusinessMemberDeleted event is fired when an additional business member is deleted
type AdditionalBusinessMemberDeleted struct {
}

func (e *AdditionalBusinessMemberDeleted) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
