package helpers

import (
	"github.com/goravel/framework/contracts/http"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/models"
)

// AuditHelper provides utilities for audit field management
type AuditHelper struct {
	permissionHelper *auth.PermissionHelper
}

// NewAuditHelper creates a new audit helper instance
func NewAuditHelper() *AuditHelper {
	return &AuditHelper{
		permissionHelper: auth.GetPermissionHelper(),
	}
}

// SetCreateAuditFields sets audit fields for creation operations
func (h *AuditHelper) SetCreateAuditFields(ctx http.Context, data map[string]interface{}) {
	user := h.permissionHelper.GetAuthenticatedUser(ctx)

	if user != nil {
		data["created_by"] = user.ID

		// Also set IP address and user agent if available
		if ipAddr := ctx.Request().Ip(); ipAddr != "" {
			data["ip_address"] = ipAddr
		}

		if userAgent := ctx.Request().Header("User-Agent"); userAgent != "" {
			data["user_agent"] = userAgent
		}
	}
}

// SetUpdateAuditFields sets audit fields for update operations
func (h *AuditHelper) SetUpdateAuditFields(ctx http.Context, data map[string]interface{}) {
	user := h.permissionHelper.GetAuthenticatedUser(ctx)

	if user != nil {
		data["updated_by"] = user.ID

		// Also set IP address and user agent if available
		if ipAddr := ctx.Request().Ip(); ipAddr != "" {
			data["ip_address"] = ipAddr
		}

		if userAgent := ctx.Request().Header("User-Agent"); userAgent != "" {
			data["user_agent"] = userAgent
		}
	}
}

// SetDeleteAuditFields sets audit fields for deletion operations
func (h *AuditHelper) SetDeleteAuditFields(ctx http.Context, data map[string]interface{}) {
	user := h.permissionHelper.GetAuthenticatedUser(ctx)

	if user != nil {
		data["deleted_by"] = user.ID

		// Also set IP address and user agent if available
		if ipAddr := ctx.Request().Ip(); ipAddr != "" {
			data["ip_address"] = ipAddr
		}

		if userAgent := ctx.Request().Header("User-Agent"); userAgent != "" {
			data["user_agent"] = userAgent
		}
	}
}

// GetAuditInfo extracts comprehensive audit information from a model
func (h *AuditHelper) GetAuditInfo(model interface{}) *models.AuditInfo {
	if auditable, ok := model.(models.AuditableInterface); ok {
		return auditable.GetAuditInfo()
	}
	return nil
}

// ValidateAuditableModel checks if a model implements the auditable interface
func (h *AuditHelper) ValidateAuditableModel(model interface{}) bool {
	_, ok := model.(models.AuditableInterface)
	return ok
}

// GetCurrentUserID returns the current authenticated user ID or nil
func (h *AuditHelper) GetCurrentUserID(ctx http.Context) *uint {
	user := h.permissionHelper.GetAuthenticatedUser(ctx)
	if user != nil {
		return &user.ID
	}
	return nil
}

// GetRequestMetadata extracts IP address and user agent from request
func (h *AuditHelper) GetRequestMetadata(ctx http.Context) (string, string) {
	ipAddr := ctx.Request().Ip()
	userAgent := ctx.Request().Header("User-Agent")
	return ipAddr, userAgent
}

// Global audit helper instance
var auditHelper *AuditHelper

// GetAuditHelper returns a singleton instance of AuditHelper
func GetAuditHelper() *AuditHelper {
	if auditHelper == nil {
		auditHelper = NewAuditHelper()
	}
	return auditHelper
}
