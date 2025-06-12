package console

import (
	"github.com/goravel/framework/contracts/console"

	"players/app/console/commands"
)

type Kernel struct {
}

// Schedule Define the application's command schedule.
// In this older signature, the scheduler instance is typically obtained via facades.Schedule() inside the method.
func (kernel *Kernel) Schedule() {
	// Example: Get the scheduler instance
	// s := facades.Schedule()

	// Example: Define a scheduled command
	// s.Command("inspire").EveryMinute()

	// If you were manually collecting events for Register:
	// var events []schedule.Event
	// events = append(events, s.Command("foo"))
	// s.Register(events)
}

// Commands Register the commands for the application.
func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.CreateAdminUser{},
	}
}
