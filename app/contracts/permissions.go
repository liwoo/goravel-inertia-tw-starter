package contracts

// StandardPermission represents the predefined permission types
type StandardPermission string

const (
	// Basic CRUD permissions
	PermissionCreate     StandardPermission = "CREATE"
	PermissionRead       StandardPermission = "READ"
	PermissionUpdate     StandardPermission = "UPDATE"
	PermissionDelete     StandardPermission = "DELETE"
	
	// Extended permissions
	PermissionExport     StandardPermission = "EXPORT"
	PermissionBulkDelete StandardPermission = "BULK_DELETE"
	PermissionBulkEdit   StandardPermission = "BULK_EDIT"
	PermissionCustom     StandardPermission = "CUSTOM"
)

// AllStandardPermissions returns all standard permission types
func AllStandardPermissions() []StandardPermission {
	return []StandardPermission{
		PermissionCreate,
		PermissionRead,
		PermissionUpdate,
		PermissionDelete,
		PermissionExport,
		PermissionBulkDelete,
		PermissionBulkEdit,
		PermissionCustom,
	}
}

// ResourcePermissionConfig defines which permissions are available for a resource
type ResourcePermissionConfig struct {
	Resource          string                        `json:"resource"`
	DisplayName       string                        `json:"display_name"`
	Category          string                        `json:"category"`
	EnabledPermissions []StandardPermission         `json:"enabled_permissions"`
	CustomPermissions  []CustomPermissionDefinition `json:"custom_permissions,omitempty"`
}

// CustomPermissionDefinition defines custom permissions beyond the standard ones
type CustomPermissionDefinition struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// PermissionRegistry holds all registered resource permissions
type PermissionRegistry struct {
	resources map[string]ResourcePermissionConfig
}

// NewPermissionRegistry creates a new permission registry
func NewPermissionRegistry() *PermissionRegistry {
	return &PermissionRegistry{
		resources: make(map[string]ResourcePermissionConfig),
	}
}

// RegisterResource registers permissions for a resource
func (r *PermissionRegistry) RegisterResource(config ResourcePermissionConfig) {
	r.resources[config.Resource] = config
}

// GetResource returns the permission config for a resource
func (r *PermissionRegistry) GetResource(resource string) (ResourcePermissionConfig, bool) {
	config, exists := r.resources[resource]
	return config, exists
}

// GetAllResources returns all registered resources
func (r *PermissionRegistry) GetAllResources() map[string]ResourcePermissionConfig {
	return r.resources
}

// GeneratePermissionSlug generates a permission slug for a resource and permission type
func GeneratePermissionSlug(resource string, permission StandardPermission) string {
	return resource + "." + string(permission)
}

// GenerateCustomPermissionSlug generates a permission slug for a custom permission
func GenerateCustomPermissionSlug(resource string, customSlug string) string {
	return resource + "." + customSlug
}