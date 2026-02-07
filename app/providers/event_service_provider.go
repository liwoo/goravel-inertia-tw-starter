package providers

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"books-database/app/events"
	"books-database/app/listeners"
)

type EventServiceProvider struct {
}

func (receiver *EventServiceProvider) Register(app foundation.Application) {
	facades.Event().Register(receiver.listen())
}

func (receiver *EventServiceProvider) Boot(app foundation.Application) {

}

func (receiver *EventServiceProvider) listen() map[event.Event][]event.Listener {
	return map[event.Event][]event.Listener{
		&events.ApplicationApproved{}: {
			&listeners.NotifyApplicationApproved{},
		},
		&events.ApplicationRejected{}: {
			&listeners.NotifyApplicationRejected{},
		},
	}
}
