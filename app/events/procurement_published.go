package events

import (
	"github.com/goravel/framework/contracts/event"
)

// ProcurementPublished is fired when a procurement notice is published
type ProcurementPublished struct {
}

func (e *ProcurementPublished) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
