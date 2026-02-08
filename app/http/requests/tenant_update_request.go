package requests

import "github.com/goravel/framework/contracts/http"

type TenantUpdateRequest struct {
	Name        *string `json:"name" form:"name"`
	Slug        *string `json:"slug" form:"slug"`
	Description *string `json:"description" form:"description"`
	IsActive    *bool   `json:"is_active" form:"is_active"`
	LogoURL     *string `json:"logo_url" form:"logo_url"`
	ID          uint    `form:"-" json:"-"` // Set by controller
}

func (r *TenantUpdateRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *TenantUpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name": "string|max:255",
		"slug": "string|max:100",
	}
}

func (r *TenantUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *TenantUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *TenantUpdateRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

func (r *TenantUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// GetResourceID returns the resource ID for update
func (r *TenantUpdateRequest) GetResourceID() interface{} {
	return r.ID
}

func (r *TenantUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}
	if r.Name != nil {
		data["name"] = *r.Name
	}
	if r.Slug != nil {
		data["slug"] = *r.Slug
	}
	if r.Description != nil {
		data["description"] = *r.Description
	}
	if r.IsActive != nil {
		data["is_active"] = *r.IsActive
	}
	if r.LogoURL != nil {
		data["logo_url"] = *r.LogoURL
	}
	return data
}
