package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/inertia"
	"players/app/services"
)

// UserPageController handles the Inertia.js User management page
// Only accessible by super admins
type UserPageController struct {
	*contracts.BasePageController
	userService *services.UserService
	authHelper  contracts.AuthHelper
}

// NewUserPageController creates a new user page controller
func NewUserPageController() *UserPageController {
	controller := &UserPageController{
		BasePageController: contracts.NewBasePageController("user", "Users/Index"),
		userService:        services.NewUserService(),
		authHelper:         helpers.NewAuthHelper(),
	}

	// Register page controller with validation
	contracts.MustRegisterPageController("users_page", controller)

	return controller
}

// checkSuperAdmin verifies if the current user is a super admin
func (c *UserPageController) checkSuperAdmin(ctx http.Context) error {
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil || !user.IsSuperAdmin {
		return fmt.Errorf("super admin access required")
	}
	return nil
}

// Index renders the Users management page with data and permissions
func (c *UserPageController) Index(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return inertia.Render(ctx, "Errors/403", map[string]interface{}{
			"message": "Access denied: Super admin privileges required",
		})
	}

	// Validate page request using contract
	req, err := c.ValidatePageRequest(ctx)
	if err != nil {
		// Return error page or redirect with error
		req = &contracts.ListRequest{Page: 1, PageSize: 20}
		req.SetDefaults()
	}

	// Build permissions map (all true for super admin)
	permissions := c.BuildPermissionsMap(ctx, "user")

	// Get users data
	usersResult, err := c.userService.GetList(*req)
	if err != nil {
		// Handle error gracefully, provide empty result
		usersResult = &contracts.PaginatedResult{
			Data:        []interface{}{},
			Total:       0,
			CurrentPage: 1,
			LastPage:    1,
			PerPage:     req.PageSize,
		}
	}

	// Get user statistics
	stats := c.getUserStatistics()

	// Get all roles for the form
	roles, _ := c.userService.GetAllRoles()

	// Build standardized page props using contract
	data := map[string]interface{}{
		"data":        usersResult.Data,
		"total":       usersResult.Total,
		"currentPage": usersResult.CurrentPage,
		"lastPage":    usersResult.LastPage,
		"perPage":     usersResult.PerPage,
		"from":        usersResult.From,
		"to":          usersResult.To,
		"hasNext":     usersResult.HasNext,
		"hasPrev":     usersResult.HasPrev,
	}

	filters := map[string]interface{}{
		"page":      req.Page,
		"pageSize":  req.PageSize,
		"search":    req.Search,
		"sort":      req.Sort,
		"direction": req.Direction,
		"filters":   req.Filters,
	}

	meta := map[string]interface{}{
		"stats": stats,
		"roles": roles,
	}

	props := c.BuildPageProps(data, filters, permissions, meta)

	return inertia.Render(ctx, "Users/Index", props)
}

// getUserStatistics returns user statistics for the dashboard
func (c *UserPageController) getUserStatistics() map[string]interface{} {
	// Get status counts
	totalCount := c.getUserTotalCount()
	activeCount := c.getUserCountByStatus(true)
	inactiveCount := totalCount - activeCount
	superAdminCount := c.getSuperAdminCount()

	return map[string]interface{}{
		"totalUsers":     totalCount,
		"activeUsers":    activeCount,
		"inactiveUsers":  inactiveCount,
		"superAdmins":    superAdminCount,
	}
}

// getUserTotalCount gets the total number of users
func (c *UserPageController) getUserTotalCount() int {
	req := contracts.ListRequest{
		PageSize: 1,
	}
	result, err := c.userService.GetList(req)
	if err != nil {
		return 0
	}
	return int(result.Total)
}

// getUserCountByStatus gets user count by active status
func (c *UserPageController) getUserCountByStatus(isActive bool) int {
	req := contracts.ListRequest{
		PageSize: 1,
		Filters: map[string]interface{}{
			"is_active": isActive,
		},
	}

	result, err := c.userService.GetListAdvanced(req, map[string]interface{}{
		"is_active": isActive,
	})
	if err != nil {
		return 0
	}

	return int(result.Total)
}

// getSuperAdminCount gets the number of super admin users
func (c *UserPageController) getSuperAdminCount() int {
	req := contracts.ListRequest{
		PageSize: 1,
		Filters: map[string]interface{}{
			"is_super_admin": true,
		},
	}

	result, err := c.userService.GetListAdvanced(req, map[string]interface{}{
		"is_super_admin": true,
	})
	if err != nil {
		return 0
	}

	return int(result.Total)
}

// CONTRACT IMPLEMENTATIONS - Required by PageControllerContract interface

// AuthorizationControllerContract implementation
func (c *UserPageController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	return c.checkSuperAdmin(ctx)
}

func (c *UserPageController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *UserPageController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *UserPageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	// For super admin only access, all permissions are based on super admin status
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	isSuperAdmin := user != nil && user.IsSuperAdmin
	
	return map[string]bool{
		"canCreate": isSuperAdmin,
		"canEdit":   isSuperAdmin,
		"canDelete": isSuperAdmin,
		"canManage": isSuperAdmin,
		"canExport": isSuperAdmin,
		"canViewReports": isSuperAdmin,
	}
}