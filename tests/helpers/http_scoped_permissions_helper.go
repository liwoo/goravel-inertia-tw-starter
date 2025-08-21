package helpers

import (
	"github.com/goravel/framework/facades"
	"players/app/models"
	"time"
)

// AssignPermissionToRole assigns a permission to a role with scope
func AssignPermissionToRole(role *models.Role, permission *models.Permission, scope string) error {
	rolePermission := &models.RolePermission{
		RoleID:       role.ID,
		PermissionID: permission.ID,
		Scope:        scope,
		GrantedAt:    time.Now(),
		IsActive:     true,
	}
	return facades.Orm().Query().Create(rolePermission)
}
