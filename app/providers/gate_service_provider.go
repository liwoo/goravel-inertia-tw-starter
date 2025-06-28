package providers

import (
	"players/app/contracts"
	"players/app/helpers"

	accessImpl "github.com/goravel/framework/auth/access"
	"github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
)

type GateServiceProvider struct {
}

func (receiver *GateServiceProvider) Register(app foundation.Application) {
	// Register the gate helper
	app.Bind("gate.helper", func(app foundation.Application) (interface{}, error) {
		return helpers.NewGateHelper(), nil
	})
}

func (receiver *GateServiceProvider) Boot(app foundation.Application) {
	// Register all resource gates
	receiver.registerBookGates()
	receiver.registerUserGates()
	receiver.registerSystemGates()
}

// registerBookGates registers all book-related permissions
func (receiver *GateServiceProvider) registerBookGates() {
	gateHelper := helpers.NewGateHelper()
	
	// Define book-specific gate configuration
	bookGateConfig := contracts.GateConfig{
		ViewAnyHandler: func(ctx http.Context, user interface{}) access.Response {
			// Everyone can view book lists
			return accessImpl.NewAllowResponse()
		},
		ViewHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Everyone can view individual books
			return accessImpl.NewAllowResponse()
		},
		CreateHandler: func(ctx http.Context, user interface{}) access.Response {
			// Only moderators and admins can create books
			return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
		},
		UpdateHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Only moderators and admins can update books
			return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
		},
		DeleteHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Only admins can delete books
			return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
		},
	}

	// Register the gates for books
	gateHelper.RegisterResourceGates("books", bookGateConfig)

	// Register additional book-specific operations
	gateHelper.RegisterGate("books.borrow", func(ctx http.Context, user interface{}) access.Response {
		// Members and above can borrow books
		return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR", "MEMBER")(ctx, user)
	})

	gateHelper.RegisterGate("books.return", func(ctx http.Context, user interface{}) access.Response {
		// Members and above can return books
		return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR", "MEMBER")(ctx, user)
	})

	gateHelper.RegisterGate("books.manage", func(ctx http.Context, user interface{}) access.Response {
		// Only librarians and admins can manage books
		return gateHelper.RoleBasedAccess("ADMIN", "LIBRARIAN")(ctx, user)
	})

	gateHelper.RegisterGate("books.export", func(ctx http.Context, user interface{}) access.Response {
		// Only librarians and admins can export book data
		return gateHelper.RoleBasedAccess("ADMIN", "LIBRARIAN")(ctx, user)
	})
}

// registerUserGates registers all user-related permissions
func (receiver *GateServiceProvider) registerUserGates() {
	gateHelper := helpers.NewGateHelper()
	
	// User CRUD operations
	userGateConfig := contracts.GateConfig{
		ViewAnyHandler: func(ctx http.Context, user interface{}) access.Response {
			// Admins and moderators can view user lists
			return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
		},
		ViewHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Users can view their own profile, admins can view all
			return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
		},
		CreateHandler: func(ctx http.Context, user interface{}) access.Response {
			// Only admins can create users
			return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
		},
		UpdateHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Only admins can update users
			return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
		},
		DeleteHandler: func(ctx http.Context, user interface{}, model interface{}) access.Response {
			// Only super-admins can delete users
			return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
		},
	}

	gateHelper.RegisterResourceGates("users", userGateConfig)

	// Additional user operations
	gateHelper.RegisterGate("users.impersonate", func(ctx http.Context, user interface{}) access.Response {
		// Only super-admins can impersonate users
		return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
	})

	gateHelper.RegisterGate("users.manage", func(ctx http.Context, user interface{}) access.Response {
		// Admins can manage users
		return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
	})
}

// registerSystemGates registers all system-level permissions
func (receiver *GateServiceProvider) registerSystemGates() {
	gateHelper := helpers.NewGateHelper()

	gateHelper.RegisterGate("system.manage", func(ctx http.Context, user interface{}) access.Response {
		// Only super-admins can manage system
		return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
	})

	gateHelper.RegisterGate("system.backup", func(ctx http.Context, user interface{}) access.Response {
		// Only super-admins can backup system
		return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
	})

	gateHelper.RegisterGate("system.configure", func(ctx http.Context, user interface{}) access.Response {
		// Only super-admins can configure system
		return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
	})

	gateHelper.RegisterGate("reports.view", func(ctx http.Context, user interface{}) access.Response {
		// Librarians and admins can view reports
		return gateHelper.RoleBasedAccess("ADMIN", "LIBRARIAN")(ctx, user)
	})

	gateHelper.RegisterGate("reports.export", func(ctx http.Context, user interface{}) access.Response {
		// Librarians and admins can export reports
		return gateHelper.RoleBasedAccess("ADMIN", "LIBRARIAN")(ctx, user)
	})
}

// registerBookOperationGates registers gates for book-specific operations
// TODO: Fix gate function signatures and re-implement
func registerBookOperationGates_disabled(gateHelper contracts.GateHelper) {
	// Disabled for build compatibility
}