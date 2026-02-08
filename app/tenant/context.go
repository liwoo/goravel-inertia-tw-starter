package tenant

import (
	"fmt"

	"books-database/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

const (
	ContextKeyTenant   = "current_tenant"
	ContextKeyTenantID = "current_tenant_id"
)

// GetFromContext retrieves the current tenant from request context
func GetFromContext(ctx http.Context) *models.Tenant {
	val := ctx.Value(ContextKeyTenant)
	if t, ok := val.(*models.Tenant); ok {
		return t
	}
	return nil
}

// GetIDFromContext returns the current tenant ID or nil
func GetIDFromContext(ctx http.Context) *uint {
	t := GetFromContext(ctx)
	if t != nil {
		return &t.ID
	}
	return nil
}

// SetInContext stores tenant in request context
func SetInContext(ctx http.Context, t *models.Tenant) {
	ctx.WithValue(ContextKeyTenant, t)
	ctx.WithValue(ContextKeyTenantID, t.ID)
}

// ResolveBySlug finds an active tenant by its URL slug
func ResolveBySlug(slug string) (*models.Tenant, error) {
	var t models.Tenant
	if err := facades.Orm().Query().
		Where("slug = ? AND is_active = ?", slug, true).
		First(&t); err != nil {
		return nil, err
	}
	if t.ID == 0 {
		return nil, fmt.Errorf("tenant not found: %s", slug)
	}
	return &t, nil
}

// GetMainTenant returns the main (default) tenant
func GetMainTenant() (*models.Tenant, error) {
	var t models.Tenant
	if err := facades.Orm().Query().
		Where("is_main = ? AND is_active = ?", true, true).
		First(&t); err != nil {
		return nil, err
	}
	if t.ID == 0 {
		return nil, fmt.Errorf("main tenant not found")
	}
	return &t, nil
}
