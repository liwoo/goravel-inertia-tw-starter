package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("tenancy", map[string]any{
		"default_slug": config.Env("DEFAULT_TENANT_SLUG", "main"),
	})
}
