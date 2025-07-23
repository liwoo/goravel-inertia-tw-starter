package auth

// PermissionScope defines the scope of access for a permission
type PermissionScope string

const (
	// ScopeByMe - User can only access resources they created
	ScopeByMe PermissionScope = "by_me"
	
	// ScopeByMyRole - User can access resources created by anyone with their role or lower
	ScopeByMyRole PermissionScope = "by_my_role"
	
	// ScopeByAll - User can access all resources regardless of creator (current default)
	ScopeByAll PermissionScope = "by_all"
)

// ScopedPermissionAction represents a permission action with a specific scope
type ScopedPermissionAction struct {
	Action CorePermissionAction `json:"action"`
	Scope  PermissionScope      `json:"scope"`
}

// GetPermissionSlug generates a permission slug for a scoped action
// Format: service_action_scope (e.g., "books_read_by_me")
func GetPermissionSlug(service ServiceRegistry, action CorePermissionAction, scope PermissionScope) string {
	return string(service) + "_" + string(action) + "_" + string(scope)
}

// GetAllPermissionScopes returns all available permission scopes
func GetAllPermissionScopes() []PermissionScope {
	return []PermissionScope{
		ScopeByMe,
		ScopeByMyRole,
		ScopeByAll,
	}
}

// GetDefaultScope returns the default permission scope
func GetDefaultScope() PermissionScope {
	return ScopeByAll
}

// IsScopedPermission checks if a permission slug includes a scope
func IsScopedPermission(slug string) bool {
	for _, scope := range GetAllPermissionScopes() {
		if len(slug) > len(scope)+1 && slug[len(slug)-len(scope)-1:] == "_"+string(scope) {
			return true
		}
	}
	return false
}

// ParsePermissionSlug extracts service, action, and scope from a permission slug
func ParsePermissionSlug(slug string) (service string, action string, scope PermissionScope) {
	// Default scope if not specified
	scope = ScopeByAll
	
	// Check if it's a scoped permission
	for _, s := range GetAllPermissionScopes() {
		suffix := "_" + string(s)
		if len(slug) > len(suffix) && slug[len(slug)-len(suffix):] == suffix {
			scope = s
			slug = slug[:len(slug)-len(suffix)]
			break
		}
	}
	
	// Split remaining slug into service and action
	for i := 0; i < len(slug); i++ {
		if slug[i] == '_' {
			service = slug[:i]
			action = slug[i+1:]
			break
		}
	}
	
	return service, action, scope
}