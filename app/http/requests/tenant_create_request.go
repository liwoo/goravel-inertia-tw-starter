package requests

import "github.com/goravel/framework/contracts/http"

type TenantCreateRequest struct {
	Name        string  `json:"name" form:"name"`
	Slug        string  `json:"slug" form:"slug"`
	Description *string `json:"description" form:"description"`
	IsActive    *bool   `json:"is_active" form:"is_active"`
	LogoURL     *string `json:"logo_url" form:"logo_url"`
}

func (r *TenantCreateRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *TenantCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name": "required|string|max:255",
		"slug": "required|string|max:100",
	}
}

func (r *TenantCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *TenantCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *TenantCreateRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

func (r *TenantCreateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

func (r *TenantCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"name": r.Name,
		"slug": r.Slug,
	}
	if r.Description != nil {
		data["description"] = *r.Description
	}
	if r.IsActive != nil {
		data["is_active"] = *r.IsActive
	} else {
		data["is_active"] = true
	}
	if r.LogoURL != nil {
		data["logo_url"] = *r.LogoURL
	}
	return data
}
