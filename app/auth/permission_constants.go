package auth

// CorePermissionActions defines the standard permission actions used across all services
type CorePermissionAction string

const (
	// Basic CRUD operations
	PermissionCreate CorePermissionAction = "create"
	PermissionRead   CorePermissionAction = "read"
	PermissionUpdate CorePermissionAction = "update"
	PermissionDelete CorePermissionAction = "delete"

	// Additional common operations
	PermissionExport     CorePermissionAction = "export"
	PermissionBulkUpdate CorePermissionAction = "bulk_update"
	PermissionBulkDelete CorePermissionAction = "bulk_delete"

	// Special operations
	PermissionManage CorePermissionAction = "manage"
	PermissionView   CorePermissionAction = "view" // For listing/viewing
)

// ServiceRegistry defines all registered services/entities in the system
type ServiceRegistry string

const (
	ServiceBooks       ServiceRegistry = "books"
	ServiceUsers       ServiceRegistry = "users"
	ServiceRoles       ServiceRegistry = "roles"
	ServicePermissions ServiceRegistry = "permissions"
	ServiceReports     ServiceRegistry = "reports"
	ServiceSystem      ServiceRegistry = "system"
	ServiceBundles     ServiceRegistry = "bundles"
)

// GetAllCorePermissionActions returns all core permission actions
func GetAllCorePermissionActions() []CorePermissionAction {
	return []CorePermissionAction{
		PermissionCreate,
		PermissionRead,
		PermissionUpdate,
		PermissionDelete,
		PermissionExport,
		PermissionBulkUpdate,
		PermissionBulkDelete,
		PermissionManage,
		PermissionView,
	}
}

// GetAllServiceRegistries returns all registered services
func GetAllServiceRegistries() []ServiceRegistry {
	return []ServiceRegistry{
		ServiceBooks,
		ServiceUsers,
		ServiceRoles,
		ServicePermissions,
		ServiceReports,
		ServiceSystem,
		ServiceBundles,
	}
}

// BuildPermissionSlug creates a permission slug in the format: service_action
func BuildPermissionSlug(service ServiceRegistry, action CorePermissionAction) string {
	return string(service) + "_" + string(action)
}

// GetServiceDisplayName returns the human-readable name for a service
func GetServiceDisplayName(service ServiceRegistry) string {
	switch service {
	case ServiceBooks:
		return "Books Management"
	case ServiceUsers:
		return "User Management"
	case ServiceRoles:
		return "Role Management"
	case ServicePermissions:
		return "Permission Management"
	case ServiceReports:
		return "Reports & Analytics"
	case ServiceSystem:
		return "System Administration"
	case ServiceBundles:
		return "SME Management"
	default:
		return string(service)
	}
}

// GetActionDisplayName returns the human-readable name for an action
func GetActionDisplayName(action CorePermissionAction) string {
	switch action {
	case PermissionCreate:
		return "Create"
	case PermissionRead:
		return "Read/List"
	case PermissionUpdate:
		return "Update/Edit"
	case PermissionDelete:
		return "Delete"
	case PermissionExport:
		return "Export"
	case PermissionBulkUpdate:
		return "Bulk Update"
	case PermissionBulkDelete:
		return "Bulk Delete"
	case PermissionManage:
		return "Full Management"
	case PermissionView:
		return "View"
	default:
		return string(action)
	}
}

// IsValidServiceAction checks if a service-action combination is valid
func IsValidServiceAction(service ServiceRegistry, action CorePermissionAction) bool {
	// All services support all actions by default
	// This can be customized per service if needed
	return true
}

// GetServiceActions returns the valid actions for a specific service
func GetServiceActions(service ServiceRegistry) []CorePermissionAction {
	switch service {
	case ServiceBooks:
		return []CorePermissionAction{
			PermissionCreate,
			PermissionRead,
			PermissionUpdate,
			PermissionDelete,
			PermissionExport,
			PermissionBulkUpdate,
			PermissionBulkDelete,
			PermissionView,
		}
	case ServiceBundles:
		return []CorePermissionAction{
			PermissionCreate,
			PermissionRead,
			PermissionUpdate,
			PermissionDelete,
			PermissionExport,
			PermissionBulkUpdate,
			PermissionBulkDelete,
			PermissionView,
		}
	case ServiceUsers:
		return []CorePermissionAction{
			PermissionCreate,
			PermissionRead,
			PermissionUpdate,
			PermissionDelete,
			PermissionExport,
			PermissionView,
			PermissionManage,
		}
	case ServiceRoles:
		return []CorePermissionAction{
			PermissionCreate,
			PermissionRead,
			PermissionUpdate,
			PermissionDelete,
			PermissionView,
			PermissionManage,
		}
	case ServicePermissions:
		return []CorePermissionAction{
			PermissionRead,
			PermissionUpdate,
			PermissionView,
			PermissionManage,
		}
	case ServiceReports:
		return []CorePermissionAction{
			PermissionView,
			PermissionExport,
		}
	case ServiceSystem:
		return []CorePermissionAction{
			PermissionView,
			PermissionManage,
		}
	default:
		return GetAllCorePermissionActions()
	}
}
