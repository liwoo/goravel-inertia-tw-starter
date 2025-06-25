package controllers

import (
	"players/app/http/inertia"

	"github.com/goravel/framework/contracts/http"
)

type UtilController struct {
	// Dependencies can be injected here
}

func NewUtilController() *UtilController {
	return &UtilController{}
}

// ShowUnaPage renders the 'Una' (Unauthorized) Inertia page.
func (r *UtilController) ShowUnaPage(ctx http.Context) http.Response {
	return inertia.Render(ctx, "Una", map[string]interface{}{
		"title": "Unauthorized Access",
	})
}
