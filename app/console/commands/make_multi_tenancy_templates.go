package commands

import "fmt"

// ──────────────────────────────────────────────
// Config template
// ──────────────────────────────────────────────

func mtTplConfig(moduleName string) string {
	return `package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("tenancy", map[string]any{
		"default_slug": config.Env("DEFAULT_TENANT_SLUG", "main"),
	})
}
`
}

// ──────────────────────────────────────────────
// Tenant context package
// ──────────────────────────────────────────────

func mtTplTenantContext(moduleName string) string {
	return fmt.Sprintf(`package tenant

import (
	"fmt"

	"%s/app/models"

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
		return nil, fmt.Errorf("tenant not found: %%s", slug)
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
`, moduleName)
}

// ──────────────────────────────────────────────
// Tenant model
// ──────────────────────────────────────────────

func mtTplTenantModel() string {
	return `package models

type Tenant struct {
	BaseAuditableModel

	Name        string  ` + "`" + `gorm:"not null;type:varchar(255)" json:"name"` + "`" + `
	Slug        string  ` + "`" + `gorm:"uniqueIndex;not null;type:varchar(100)" json:"slug"` + "`" + `
	IsActive    bool    ` + "`" + `gorm:"default:true;index" json:"is_active"` + "`" + `
	IsMain      bool    ` + "`" + `gorm:"default:false;index" json:"is_main"` + "`" + `
	Description *string ` + "`" + `gorm:"type:text" json:"description,omitempty"` + "`" + `
	Settings    *string ` + "`" + `gorm:"type:jsonb" json:"settings,omitempty"` + "`" + `
	LogoURL     *string ` + "`" + `gorm:"type:varchar(500)" json:"logo_url,omitempty"` + "`" + `
}

func (Tenant) TableName() string {
	return "tenants"
}
`
}

// ──────────────────────────────────────────────
// UserTenant model
// ──────────────────────────────────────────────

func mtTplUserTenantModel() string {
	return `package models

import "time"

type UserTenant struct {
	BaseAuditableModel

	UserID   uint      ` + "`" + `gorm:"not null;index;uniqueIndex:idx_user_tenant" json:"user_id"` + "`" + `
	TenantID uint      ` + "`" + `gorm:"not null;index;uniqueIndex:idx_user_tenant" json:"tenant_id"` + "`" + `
	RoleID   *uint     ` + "`" + `gorm:"index" json:"role_id,omitempty"` + "`" + `
	IsActive bool      ` + "`" + `gorm:"default:true" json:"is_active"` + "`" + `
	JoinedAt time.Time ` + "`" + `gorm:"not null;autoCreateTime" json:"joined_at"` + "`" + `

	// Relations
	User   *User   ` + "`" + `gorm:"foreignKey:UserID" json:"user,omitempty"` + "`" + `
	Tenant *Tenant ` + "`" + `gorm:"foreignKey:TenantID" json:"tenant,omitempty"` + "`" + `
	Role   *Role   ` + "`" + `gorm:"foreignKey:RoleID" json:"role,omitempty"` + "`" + `
}

func (UserTenant) TableName() string {
	return "user_tenants"
}
`
}

// ──────────────────────────────────────────────
// Migrations
// ──────────────────────────────────────────────

func mtTplMigrationTenants(timestamp string) string {
	structName := fmt.Sprintf("M%sCreateTenantsTable", timestamp)
	sig := fmt.Sprintf("%s_create_tenants_table", timestamp)
	return fmt.Sprintf(`package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type %s struct{}

func (r *%s) Signature() string {
	return "%s"
}

func (r *%s) Up() error {
	if !facades.Schema().HasTable("tenants") {
		return facades.Schema().Create("tenants", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 255)
			table.String("slug", 100)
			table.Index("slug")
			table.Boolean("is_active").Default(true)
			table.Boolean("is_main").Default(false)
			table.Text("description").Nullable()
			table.Text("settings").Nullable()
			table.String("logo_url", 500).Nullable()

			// Audit fields
			table.UnsignedBigInteger("tenant_id").Nullable()
			table.UnsignedBigInteger("created_by").Nullable()
			table.UnsignedBigInteger("updated_by").Nullable()
			table.UnsignedBigInteger("deleted_by").Nullable()
			table.String("ip_address", 45).Nullable()
			table.Text("user_agent").Nullable()

			table.TimestampsTz()
			table.SoftDeletesTz()
		})
	}
	return nil
}

func (r *%s) Down() error {
	return facades.Schema().DropIfExists("tenants")
}
`, structName, structName, sig, structName, structName)
}

func mtTplMigrationUserTenants(timestamp string) string {
	structName := fmt.Sprintf("M%sCreateUserTenantsTable", timestamp)
	sig := fmt.Sprintf("%s_create_user_tenants_table", timestamp)
	return fmt.Sprintf(`package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type %s struct{}

func (r *%s) Signature() string {
	return "%s"
}

func (r *%s) Up() error {
	if !facades.Schema().HasTable("user_tenants") {
		return facades.Schema().Create("user_tenants", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.UnsignedBigInteger("tenant_id")
			table.UnsignedBigInteger("role_id").Nullable()
			table.Boolean("is_active").Default(true)
			table.TimestampTz("joined_at").UseCurrent()

			// Audit fields
			table.UnsignedBigInteger("tenant_id_ref").Nullable()
			table.UnsignedBigInteger("created_by").Nullable()
			table.UnsignedBigInteger("updated_by").Nullable()
			table.UnsignedBigInteger("deleted_by").Nullable()
			table.String("ip_address", 45).Nullable()
			table.Text("user_agent").Nullable()

			table.TimestampsTz()
			table.SoftDeletesTz()

			// Indexes and constraints
			table.Index("user_id")
			table.Index("tenant_id")
		})
	}
	return nil
}

func (r *%s) Down() error {
	return facades.Schema().DropIfExists("user_tenants")
}
`, structName, structName, sig, structName, structName)
}

func mtTplMigrationAddTenantID(timestamp string) string {
	structName := fmt.Sprintf("M%sAddTenantIdToAllTables", timestamp)
	sig := fmt.Sprintf("%s_add_tenant_id_to_all_tables", timestamp)
	return fmt.Sprintf(`package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type %s struct{}

func (r *%s) Signature() string {
	return "%s"
}

func (r *%s) Up() error {
	tables := []string{
		"users", "books", "lenders", "applications",
		"configs", "messages", "notifications", "user_activities",
		"authors",
	}

	for _, tableName := range tables {
		if facades.Schema().HasTable(tableName) && !facades.Schema().HasColumn(tableName, "tenant_id") {
			if err := facades.Schema().Table(tableName, func(table schema.Blueprint) {
				table.UnsignedBigInteger("tenant_id").Nullable()
				table.Index("tenant_id")
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *%s) Down() error {
	tables := []string{
		"users", "books", "lenders", "applications",
		"configs", "messages", "notifications", "user_activities",
		"authors",
	}

	for _, tableName := range tables {
		if facades.Schema().HasTable(tableName) && facades.Schema().HasColumn(tableName, "tenant_id") {
			facades.Schema().Table(tableName, func(table schema.Blueprint) {
				table.DropColumn("tenant_id")
			})
		}
	}

	return nil
}
`, structName, structName, sig, structName, structName)
}

// ──────────────────────────────────────────────
// Tenant seeder
// ──────────────────────────────────────────────

func mtTplTenantSeeder(moduleName string) string {
	return fmt.Sprintf(`package seeders

import (
	"%s/app/models"

	"github.com/goravel/framework/facades"
)

type TenantSeeder struct{}

func (s *TenantSeeder) Signature() string {
	return "TenantSeeder"
}

func (s *TenantSeeder) Run() error {
	// Check if main tenant already exists
	var existing models.Tenant
	facades.Orm().Query().Where("is_main = ?", true).First(&existing)
	if existing.ID > 0 {
		facades.Log().Info("Main tenant already exists, skipping...")
		return nil
	}

	// Get default slug from config
	slug := facades.Config().GetString("tenancy.default_slug", "main")

	mainTenant := models.Tenant{
		Name:     facades.Config().GetString("app.name", "Main"),
		Slug:     slug,
		IsActive: true,
		IsMain:   true,
	}

	if err := facades.Orm().Query().Create(&mainTenant); err != nil {
		return err
	}
	facades.Log().Info("Main tenant created successfully")

	// Assign all existing super admin users to the main tenant
	var superAdmins []models.User
	facades.Orm().Query().Where("is_super_admin = ?", true).Find(&superAdmins)

	for _, admin := range superAdmins {
		ut := models.UserTenant{
			UserID:   admin.ID,
			TenantID: mainTenant.ID,
			IsActive: true,
		}
		if err := facades.Orm().Query().Create(&ut); err != nil {
			facades.Log().Warning("Failed to assign super admin to main tenant: " + err.Error())
		}
	}

	// Assign all other existing users to the main tenant
	var users []models.User
	facades.Orm().Query().Where("is_super_admin = ? OR is_super_admin IS NULL", false).Find(&users)

	for _, user := range users {
		ut := models.UserTenant{
			UserID:   user.ID,
			TenantID: mainTenant.ID,
			IsActive: true,
		}
		facades.Orm().Query().Create(&ut)
	}

	facades.Log().Info("All existing users assigned to main tenant")
	return nil
}
`, moduleName)
}

// ──────────────────────────────────────────────
// Tenant middleware
// ──────────────────────────────────────────────

func mtTplTenantMiddleware(moduleName string) string {
	return fmt.Sprintf(`package middleware

import (
	"%s/app/tenant"

	contractshttp "github.com/goravel/framework/contracts/http"
)

// TenantFromSlug resolves the current tenant from the {tenant} route parameter
func TenantFromSlug() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		slug := ctx.Request().Route("tenant")

		if slug == "" {
			// No tenant slug in URL - resolve main tenant
			mainTenant, err := tenant.GetMainTenant()
			if err != nil {
				ctx.Request().AbortWithStatusJson(500, contractshttp.Json{
					"error": "Failed to resolve main tenant",
				})
				return
			}
			tenant.SetInContext(ctx, mainTenant)
			ctx.Request().Next()
			return
		}

		// Resolve tenant from slug
		t, err := tenant.ResolveBySlug(slug)
		if err != nil {
			ctx.Request().AbortWithStatusJson(404, contractshttp.Json{
				"error": "Tenant not found",
			})
			return
		}

		tenant.SetInContext(ctx, t)
		ctx.Request().Next()
	}
}
`, moduleName)
}

// ──────────────────────────────────────────────
// Tenant service
// ──────────────────────────────────────────────

func mtTplTenantService(moduleName string) string {
	return fmt.Sprintf(`package services

import (
	"fmt"

	"%s/app/contracts"
	"%s/app/models"

	"github.com/goravel/framework/facades"
)

type TenantService struct {
	contracts.CrudServiceContract
}

func NewTenantService() *TenantService {
	service := contracts.NewServiceBuilder[models.Tenant]("tenants", "id").
		WithSearchFields("name", "slug", "description").
		WithSortFields("id", "name", "slug", "created_at", "updated_at").
		WithFilterFields("is_active", "is_main").
		WithValidationRules(map[string]interface{}{
			"name": "required|string|max:255",
			"slug": "required|string|max:100",
		}).
		WithRelations("Creator", "Updater").
		WithDefaultSort("created_at", "DESC").
		WithSoftDeletes().
		Build()

	svc := &TenantService{CrudServiceContract: service}
	contracts.SetActualServiceHelper(service, svc, "TenantService")
	return svc
}

// AssignUserToTenant adds a user to a tenant
func (s *TenantService) AssignUserToTenant(userID, tenantID uint, roleID *uint) error {
	// Check if already assigned
	var existing models.UserTenant
	facades.Orm().Query().
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		First(&existing)
	if existing.ID > 0 {
		// Reactivate if inactive
		if !existing.IsActive {
			_, err := facades.Orm().Query().Model(&existing).Update("is_active", true)
			return err
		}
		return fmt.Errorf("user is already assigned to this tenant")
	}

	ut := models.UserTenant{
		UserID:   userID,
		TenantID: tenantID,
		RoleID:   roleID,
		IsActive: true,
	}
	return facades.Orm().Query().Create(&ut)
}

// RemoveUserFromTenant deactivates a user's tenant membership
func (s *TenantService) RemoveUserFromTenant(userID, tenantID uint) error {
	_, err := facades.Orm().Query().
		Model(&models.UserTenant{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Update("is_active", false)
	return err
}

// GetTenantUsers returns users assigned to a tenant
func (s *TenantService) GetTenantUsers(tenantID uint) ([]models.UserTenant, error) {
	var userTenants []models.UserTenant
	if err := facades.Orm().Query().
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		With("User").
		With("Role").
		Find(&userTenants); err != nil {
		return nil, err
	}
	return userTenants, nil
}
`, moduleName, moduleName)
}

// ──────────────────────────────────────────────
// Tenant controller
// ──────────────────────────────────────────────

func mtTplTenantController(moduleName string) string {
	return fmt.Sprintf(`package tenants

import (
	"fmt"
	"strconv"

	"%s/app/auth"
	"%s/app/contracts"
	"%s/app/http/requests"
	"%s/app/models"
	"%s/app/services"

	"github.com/goravel/framework/contracts/http"
)

type TenantController struct {
	*contracts.CrudController[models.Tenant, *requests.TenantCreateRequest, *requests.TenantUpdateRequest]
	tenantService *services.TenantService
}

func NewTenantController() *TenantController {
	tenantService := services.NewTenantService()

	crudController := contracts.NewCrudController[models.Tenant, *requests.TenantCreateRequest, *requests.TenantUpdateRequest]("tenant", tenantService).
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			permHelper := auth.GetPermissionHelper()
			user, err := permHelper.RequireAuthentication(ctx)
			if err != nil {
				return err
			}
			if !user.IsSuperAdminUser() {
				return fmt.Errorf("super admin access required")
			}
			return nil
		}).
		Build()

	return &TenantController{
		CrudController: crudController,
		tenantService:  tenantService,
	}
}

func (c *TenantController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

func (c *TenantController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

func (c *TenantController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

func (c *TenantController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

func (c *TenantController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// AssignUser assigns a user to the tenant
func (c *TenantController) AssignUser(ctx http.Context) http.Response {
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	tenantID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid tenant ID", nil)
	}

	userIDStr := ctx.Request().Input("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	var roleID *uint
	if roleIDStr := ctx.Request().Input("role_id"); roleIDStr != "" {
		if id, err := strconv.ParseUint(roleIDStr, 10, 64); err == nil {
			uid := uint(id)
			roleID = &uid
		}
	}

	if err := c.tenantService.AssignUserToTenant(uint(userID), tenantID, roleID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "User assigned to tenant")
}

// RemoveUser removes a user from the tenant
func (c *TenantController) RemoveUser(ctx http.Context) http.Response {
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	tenantID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid tenant ID", nil)
	}

	userIDStr := ctx.Request().Route("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	if err := c.tenantService.RemoveUserFromTenant(uint(userID), tenantID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "User removed from tenant")
}

// GetUsers returns users assigned to the tenant
func (c *TenantController) GetUsers(ctx http.Context) http.Response {
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	tenantID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid tenant ID", nil)
	}

	users, err := c.tenantService.GetTenantUsers(tenantID)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve tenant users: "+err.Error())
	}

	return c.SuccessResponse(ctx, users, "Tenant users retrieved")
}
`, moduleName, moduleName, moduleName, moduleName, moduleName)
}

// ──────────────────────────────────────────────
// Tenant page controller
// ──────────────────────────────────────────────

func mtTplTenantPageController(moduleName string) string {
	return fmt.Sprintf(`package tenants

import (
	"%s/app/auth"
	"%s/app/contracts"
	"%s/app/services"
)

type TenantPageController struct {
	*contracts.GenericPageController
	tenantService *services.TenantService
}

func NewTenantPageController() *TenantPageController {
	tenantService := services.NewTenantService()
	return &TenantPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "tenants",
			PageComponent:     "Tenant/Index",
			Service:           tenantService,
			ServiceIdentifier: auth.ServiceTenants,
			StatsEnabled:      false,
		}),
		tenantService: tenantService,
	}
}
`, moduleName, moduleName, moduleName)
}

// ──────────────────────────────────────────────
// Tenant request files
// ──────────────────────────────────────────────

func mtTplTenantCreateRequest() string {
	bt := "`"
	return fmt.Sprintf(`package requests

import "github.com/goravel/framework/contracts/http"

type TenantCreateRequest struct {
	Name        string  %sjson:"name" form:"name"%s
	Slug        string  %sjson:"slug" form:"slug"%s
	Description *string %sjson:"description" form:"description"%s
	IsActive    *bool   %sjson:"is_active" form:"is_active"%s
	LogoURL     *string %sjson:"logo_url" form:"logo_url"%s
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
`, bt, bt, bt, bt, bt, bt, bt, bt, bt, bt)
}

func mtTplTenantUpdateRequest() string {
	bt := "`"
	return fmt.Sprintf(`package requests

import "github.com/goravel/framework/contracts/http"

type TenantUpdateRequest struct {
	Name        *string %sjson:"name" form:"name"%s
	Slug        *string %sjson:"slug" form:"slug"%s
	Description *string %sjson:"description" form:"description"%s
	IsActive    *bool   %sjson:"is_active" form:"is_active"%s
	LogoURL     *string %sjson:"logo_url" form:"logo_url"%s
	ID          uint    %sform:"-" json:"-"%s // Set by controller
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
`, bt, bt, bt, bt, bt, bt, bt, bt, bt, bt, bt, bt)
}
