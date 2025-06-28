package contracts

import (
	"github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/http"
)

// Standard CRUD permissions for any resource
const (
	ViewAnyPermission     = "viewAny"     // List resources
	ViewPermission        = "view"        // View specific resource
	CreatePermission      = "create"      // Create new resource
	UpdatePermission      = "update"      // Update existing resource
	DeletePermission      = "delete"      // Delete resource
	RestorePermission     = "restore"     // Restore soft-deleted resource
	ForceDeletePermission = "forceDelete" // Permanently delete
)

// GateConfig defines handlers for different permissions
type GateConfig struct {
	ViewAnyHandler func(ctx http.Context, user interface{}) access.Response
	ViewHandler    func(ctx http.Context, user interface{}, model interface{}) access.Response
	CreateHandler  func(ctx http.Context, user interface{}) access.Response
	UpdateHandler  func(ctx http.Context, user interface{}, model interface{}) access.Response
	DeleteHandler  func(ctx http.Context, user interface{}, model interface{}) access.Response
}

// AuthHelper provides authorization helper functions
type AuthHelper interface {
	// Role checks
	HasRole(user interface{}, roles ...string) bool
	HasAnyRole(user interface{}, roles []string) bool
	HasAllRoles(user interface{}, roles []string) bool

	// Permission checks
	HasPermission(user interface{}, permission string) bool

	// Ownership checks
	IsOwner(user interface{}, resource interface{}) bool

	// Context checks
	IsSuperAdmin(user interface{}) bool
	IsAuthenticated(ctx http.Context) bool
	GetCurrentUser(ctx http.Context) interface{}

	// Resource-specific checks
	CanManageResource(user interface{}, resource string) bool
	CanAccessResource(user interface{}, resource interface{}) bool
}

// GateHelper provides common authorization patterns
type GateHelper interface {
	// Common gate patterns
	RoleBasedAccess(allowedRoles ...string) func(ctx http.Context, user interface{}) access.Response
	OwnershipBasedAccess() func(ctx http.Context, user interface{}, model interface{}) access.Response
	ConditionalAccess(condition func(ctx http.Context, user interface{}) bool, message string) func(ctx http.Context, user interface{}) access.Response

	// Gate registration
	RegisterResourceGates(resource string, config GateConfig)
	RegisterGate(name string, handler func(ctx http.Context, user interface{}) access.Response)
}