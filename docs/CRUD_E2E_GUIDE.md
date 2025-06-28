# Complete CRUD Service Implementation Guide

This guide walks you through creating a complete, production-ready CRUD system from scratch using our automated scaffolding tools. Perfect for junior developers who want to understand how all the pieces fit together.

## 🚀 Quick Start (TL;DR)

```bash
# 1. Generate complete CRUD system
go run . artisan make:crud-e2e Product

# 2. Run migrations
go run . artisan migrate

# 3. Seed permissions
go run . artisan seed --seeder=rbac

# 4. Start developing!
```

## 📋 Prerequisites

Before starting, ensure you have:
- ✅ RBAC system set up: `go run . artisan rbac:setup`
- ✅ Admin user created: `go run . artisan user:create-admin`
- ✅ Frontend build process running: `npm run dev`
- ✅ Fresh database migrations: `go run . artisan migrate:fresh`
- ✅ Seeded permissions: `go run . artisan db:seed --seeder=DatabaseSeeder`

**Note:** If you're having permission issues, run fresh migrations and re-seed the database to ensure all schema fixes are applied.

**Service Restart:** After making backend code changes (controllers, services, models), you'll need to restart your Go server to see the changes:
- Stop the server with `Ctrl+C` (or `Cmd+C` on Mac)
- Restart with `go run .` or `air` (if using air for hot reload)

## 🎯 What You'll Build

Our `make:crud-e2e` command generates a **complete, production-ready CRUD system** with:

### Backend Components (7 files)
1. **Model** with soft deletes and validation
2. **Migration** with proper indexes and foreign keys
3. **Service** implementing all CRUD contracts
4. **Request Validation** classes for create/update
5. **API Controller** with permission enforcement
6. **Page Controller** for Inertia.js server-side rendering
7. **Routes** with proper middleware and permissions

### Frontend Components (4 file groups)
8. **TypeScript Types** for full type safety
9. **React Components** (columns, forms, detail views)
10. **React Pages** with complete CRUD functionality
11. **Permission Integration** for role-based UI

---

## 📖 Step-by-Step Implementation

### Step 1: Generate Your CRUD System

```bash
# Replace "Product" with your resource name (singular, PascalCase)
go run . artisan make:crud-e2e Product
```

**⚠️ Important Naming Convention:**
- Use **singular** form for the command (e.g., `Product`, not `Products`)
- The system will automatically pluralize for:
  - Database table names (`products`)
  - Route paths (`/admin/products`)
  - Permission slugs (`products.create`, `products.view`, etc.)
  - API endpoints (`/api/products`)
- This ensures consistency with the RBAC permission system

**What this generates:**
```
🔨 Creating model...
✓ Creating model generated successfully
🔨 Creating migration...
✓ Creating migration generated successfully
🔨 Creating service with contracts...
✓ Creating service with contracts generated successfully
🔨 Creating validation requests...
✓ Creating validation requests generated successfully
🔨 Creating API controller...
✓ Creating API controller generated successfully
🔨 Creating page controller...
✓ Creating page controller generated successfully
🔨 Adding routes...
✓ Adding routes generated successfully
🔨 Creating permissions...
✓ Creating permissions generated successfully
🔨 Creating TypeScript types...
✓ Creating TypeScript types generated successfully
🔨 Creating React components...
✓ Creating React components generated successfully
🔨 Creating React pages...
✓ Creating React pages generated successfully

🎉 Complete CRUD system generated successfully!
```

### Step 2: Review Generated Files

#### Backend Files Structure
```
app/
├── models/product.go                    # Model with soft deletes
├── services/product_service.go          # Service with contracts
├── http/
│   ├── controllers/
│   │   ├── product_controller.go        # API controller
│   │   └── products_page_controller.go  # Page controller
│   └── requests/
│       ├── create_product_request.go    # Create validation
│       └── update_product_request.go    # Update validation
└── database/
    ├── migrations/
    │   └── *_create_products_table.go   # Migration with indexes
    └── seeders/
        └── product_permission_seeder.go # RBAC permissions

routes/products.go                       # Route definitions
```

#### Frontend Files Structure
```
resources/js/
├── types/product.ts                     # TypeScript definitions
├── components/Products/
│   ├── ProductColumns.tsx              # Table columns definition
│   ├── ProductForm.tsx                 # Create/Edit form
│   └── ProductDetail.tsx               # Detail view for modals
└── pages/Products/
    └── Index.tsx                        # Main CRUD page
```

### Step 3: Run Database Migration

```bash
# Apply the new migration
go run . artisan migrate
```

This creates your table with:
- Primary key (`id`)
- All your specified fields
- Soft delete support (`deleted_at`)
- Timestamps (`created_at`, `updated_at`)
- Proper indexes for performance

### Step 4: Seed Permissions

```bash
# Seed the RBAC permissions for your new resource
go run . artisan seed --seeder=rbac
```

This creates permissions like:
- `products.view` - View product listings
- `products.create` - Create new products
- `products.edit` - Edit existing products
- `products.delete` - Delete products

### Step 5: Update Frontend Routing

Add your new page to the frontend router:

```typescript
// In your main router file
import ProductsIndex from '@/pages/Products/Index';

// Add route
<Route path="/admin/products" component={ProductsIndex} />
```

### Step 6: Test Your CRUD System

1. **Visit the page**: Navigate to `/admin/products`
2. **Test permissions**: Try with different user roles
3. **Test CRUD operations**:
   - ✅ Create new records
   - ✅ View and search/filter listings
   - ✅ Edit existing records
   - ✅ Delete records (soft delete)
   - ✅ Pagination and sorting

---

## 🔧 Understanding the Generated Code

### 1. Model (`app/models/product.go`)

```go
type Product struct {
    ID          uint      `json:"id" gorm:"primarykey"`
    Name        string    `json:"name" gorm:"not null;size:255;index"`
    Description string    `json:"description" gorm:"type:text"`
    Price       float64   `json:"price" gorm:"not null;index"`
    IsActive    bool      `json:"is_active" gorm:"default:true;index"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
```

**Key Features:**
- ✅ Soft deletes with `DeletedAt`
- ✅ Proper GORM tags for database constraints
- ✅ JSON tags for API serialization
- ✅ Indexes on searchable/filterable fields

### 2. Service Contract Implementation

```go
type ProductService struct {
    authHelper contracts.AuthHelperContract
}

// Implements ALL required contracts:
// - CrudServiceContract
// - PaginationServiceContract  
// - SortableServiceContract
// - FilterableServiceContract
```

**Contract Enforcement:**
- ✅ Can't compile without implementing all methods
- ✅ Standardized pagination, sorting, filtering
- ✅ Type-safe error handling
- ✅ Permission integration ready

### 3. Controller with Permission Enforcement

```go
func (c *ProductController) Store(ctx http.Context) http.Response {
    // Automatic permission check
    if err := c.CheckPermission(ctx, "products.create", nil); err != nil {
        return c.ForbiddenResponse(ctx, "Access denied")
    }
    
    // Validation using generated request class
    var req requests.CreateProductRequest
    if err := ctx.Request().Bind(&req); err != nil {
        return c.ValidationErrorResponse(ctx, err)
    }
    
    // Service call with error handling
    product, err := c.productService.Create(&req)
    if err != nil {
        return c.ErrorResponse(ctx, "Failed to create product", err)
    }
    
    return c.SuccessResponse(ctx, "Product created successfully", product)
}
```

**Built-in Features:**
- ✅ Permission checking on every action
- ✅ Request validation with custom classes
- ✅ Standardized error responses
- ✅ Service layer integration

### 4. Frontend with Type Safety

```typescript
// Generated TypeScript types
export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

// Form validation
const formSchema = z.object({
  name: z.string().min(1, "Name is required").max(255),
  description: z.string().optional(),
  price: z.number().min(0, "Price must be positive"),
  is_active: z.boolean().default(true),
});
```

**Type Safety Benefits:**
- ✅ Compile-time error checking
- ✅ IDE autocompletion
- ✅ Runtime validation with Zod
- ✅ Consistent data structures

---

## 🎨 Customization Guide

### Adding Custom Fields

1. **Update the model**:
```go
type Product struct {
    // ... existing fields
    Category    string    `json:"category" gorm:"size:100;index"`
    SKU         string    `json:"sku" gorm:"unique;size:50"`
}
```

2. **Update the migration**:
```go
err = facades.Schema().Create("products", func(table schema.Blueprint) {
    // ... existing fields
    table.String("category", 100).Index()
    table.String("sku", 50).Unique()
})
```

3. **Update TypeScript types**:
```typescript
export interface Product {
  // ... existing fields
  category: string;
  sku: string;
}
```

4. **Update form validation**:
```typescript
const formSchema = z.object({
  // ... existing fields
  category: z.string().min(1).max(100),
  sku: z.string().min(1).max(50),
});
```

### Adding Custom Business Logic

**Service Layer** (`app/services/product_service.go`):
```go
func (s *ProductService) Create(req *requests.CreateProductRequest) (*models.Product, error) {
    // Custom business logic before creation
    if err := s.validateSKUUnique(req.SKU); err != nil {
        return nil, err
    }
    
    // Call parent implementation
    return s.BaseCrudService.Create(req)
    
    // Custom logic after creation (notifications, etc.)
}
```

### Adding Relationships

```go
type Product struct {
    // ... existing fields
    CategoryID uint     `json:"category_id" gorm:"index"`
    Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
}
```

---

## 🛡️ Security & Best Practices

### Permission System with Auto-Detection

The permission system now features automatic detection and server-side enforcement:

#### Page Controller (Server-Side Enforcement)
```go
// Every page controller MUST check permissions before rendering
func (c *ProductsPageController) Index(ctx http.Context) http.Response {
    // Server-side permission check - returns 403 if unauthorized
    permHelper := auth.GetPermissionHelper()
    _, err := permHelper.RequireServicePermission(ctx, auth.ServiceProducts, auth.PermissionRead)
    if err != nil {
        return ctx.Response().Status(403).Json(map[string]interface{}{
            "error": "Forbidden",
            "message": "You don't have permission to access this page",
        })
    }
    
    // Permissions are automatically included in global props
    return inertia.Render(ctx, "Products/Index", props)
}
```

#### Frontend (Auto-Detection)
```tsx
// No manual permission props needed!
<CrudPage
    resourceName="products"  // Automatically detects all permissions
    title="Products Management"
    columns={productColumns}
    data={data}
    filters={filters}
    // No canCreate, canEdit, canDelete props needed!
/>
```

#### Navigation (Auto-Filtering)
```tsx
// Configure navigation with permission requirements
const navigationConfig = {
    navMain: [
        {
            title: "Products",
            url: "/admin/products",
            icon: PackageIcon,
            requiredService: "products",
            requiredAction: "read" as const,  // Menu item only shows if user has products_read
        },
    ]
}
```

### Input Validation

All requests use dedicated validation classes:

```go
type CreateProductRequest struct {
    Name        string  `json:"name" form:"name" validate:"required,max=255"`
    Description string  `json:"description" form:"description"`
    Price       float64 `json:"price" form:"price" validate:"required,min=0"`
    IsActive    bool    `json:"is_active" form:"is_active"`
}
```

**⚠️ Critical Validation Issue & Fix**

The generated validation system has a known issue with Goravel's `ValidateRequest()` method that can cause "unexpected end of JSON input" errors and incorrect validation failures. Here's the **required fix** for your validation methods:

#### Problem
The default generated validation method tries to use `ValidateRequest()` which has issues with request body consumption:

```go
// ❌ Problematic - causes EOF errors
func (c *ProductController) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
    var createRequest requests.ProductCreateRequest
    errors, err := ctx.Request().ValidateRequest(&createRequest)  // This fails!
    if err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    // ...
}
```

#### Solution
**Replace** your validation methods with manual binding and validation:

```go
// ✅ Working solution - use manual binding
func (c *ProductController) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
    var createRequest requests.ProductCreateRequest
    
    // Bind the data to the struct
    if err := ctx.Request().Bind(&createRequest); err != nil {
        return nil, fmt.Errorf("data binding failed: %w", err)
    }
    
    // Manual validation - check field lengths
    if len(createRequest.Name) > 255 {
        return nil, fmt.Errorf("validation errors: name exceeds 255 characters (%d)", len(createRequest.Name))
    }
    if createRequest.Price < 0 {
        return nil, fmt.Errorf("validation errors: price must be positive")
    }
    
    // Check required fields
    if createRequest.Name == "" {
        return nil, fmt.Errorf("validation errors: name is required")
    }
    
    return createRequest.ToCreateData(), nil
}
```

#### Why This Happens
1. `ValidateRequest()` tries to read the request body for validation
2. The request body stream can only be read once in HTTP
3. If the body was consumed elsewhere, validation gets empty data
4. Empty strings incorrectly trigger max length validation errors

#### Quick Fix for Existing Projects
If you have existing validation issues, update your controller validation methods:

1. **Remove** `ctx.Request().ValidateRequest(&request)` calls
2. **Add** `ctx.Request().Bind(&request)` instead  
3. **Implement** manual validation rules
4. **Test** that your validation now works correctly

This fix ensures reliable validation for all CRUD operations.

### SQL Injection Prevention

GORM ORM with parameterized queries:

```go
// Safe - uses parameterized queries
err := facades.Orm().Query().Where("name LIKE ?", "%"+search+"%").Find(&products)
```

---

## 🧪 Testing Your Implementation

### 1. Backend API Testing

```bash
# Test API endpoints
curl -X GET "http://localhost:3500/api/products"
curl -X POST "http://localhost:3500/api/products" -d '{"name":"Test Product","price":29.99}'
```

### 2. Permission Testing

```bash
# Create users with different roles
go run . artisan user:create test@example.com password123
go run . artisan role:assign test@example.com member

# Test permission matrix
# 1. Visit /admin/permissions as super admin
# 2. Assign specific permissions to roles
# 3. Login as test user to verify access

# Test server-side enforcement
curl -X GET "http://localhost:3500/admin/products" -H "Cookie: your-session-cookie"
# Should return 403 if no products_read permission
```

#### Using Permission Hooks in Custom Components

```tsx
import { usePermissions } from '@/contexts/PermissionsContext';

function MyProductComponent() {
    const { canPerformAction, isSuperAdmin } = usePermissions();
    
    // Check specific permission
    if (canPerformAction('products', 'create')) {
        // Show create button
    }
    
    // Check super admin status
    if (isSuperAdmin()) {
        // Show admin features
    }
}
```

### 3. Frontend Testing

1. **Create Operations**: Test form validation and success states
2. **Read Operations**: Test search, pagination, sorting
3. **Update Operations**: Test inline editing and validation
4. **Delete Operations**: Test soft delete confirmation

### 4. Debugging Permissions

To debug permission issues, add temporary logging:

```go
// In your page controller
permissions := c.BuildPermissionsMap(ctx, "books")
fmt.Printf("DEBUG: Permissions for user: %+v\n", permissions)

// In permission helper
fmt.Printf("DEBUG: User %d has %d roles loaded\n", user.ID, len(user.Roles))
for _, role := range user.Roles {
    fmt.Printf("DEBUG: User has role: %s\n", role.Slug)
}

// In permission service
if user.IsSuperAdmin() {
    fmt.Printf("DEBUG: User %d is super admin\n", user.ID)
}
```

Check the server logs to see:
- Whether the user is authenticated
- What roles are loaded
- What permissions are granted
- Whether super-admin status is recognized

---

## 🚨 Troubleshooting

### Common Issues

**Migration Errors:**
```bash
# Drop and recreate if needed
go run . artisan migrate:rollback
go run . artisan migrate
```

**Permission Denied:**
```bash
# Re-seed permissions
go run . artisan seed --seeder=rbac
```

**Type Errors in Frontend:**
```bash
# Regenerate types after model changes
npm run build
```

**Service Contract Errors:**
```go
// Ensure your service implements ALL required methods
type ProductService struct {
    // Must implement:
    // - GetList, GetListAdvanced
    // - Create, Update, Delete, GetByID
    // - GetPaginated, GetTotalCount
    // - GetSorted, GetSortOptions
    // - GetFiltered, GetFilterOptions
}
```

**Validation Errors (EOF, "unexpected end of JSON input"):**
```go
// ❌ Problem: ValidateRequest() consuming request body
errors, err := ctx.Request().ValidateRequest(&request)

// ✅ Solution: Use manual binding instead
if err := ctx.Request().Bind(&request); err != nil {
    return nil, fmt.Errorf("data binding failed: %w", err)
}
// Add manual validation rules here
```

**Development Workflow Issues:**

1. **Backend Changes Not Reflecting:**
   ```bash
   # Always restart the Go server after backend changes
   # Press Ctrl+C (or Cmd+C on Mac) to stop
   go run .
   # Or if using air for hot reload:
   air
   ```

2. **Frontend Changes Not Loading:**
   ```bash
   # Ensure your frontend build process is running
   npm run dev
   # Or yarn dev
   ```

### Permission System Issues

**Issue: Permissions not reflecting in UI**

The permission system now loads permissions globally and components auto-detect them:

1. **Debug Permission Loading:**
```go
// Check console for debug output
DEBUG HasPermission: user 1 has permissions: [products_create, products_read]
DEBUG HasPermission: checking permission: products_update
DEBUG loadUserPermissions: role member has 2 permissions
```

2. **Verify Permission Format:**
```go
// ✅ Correct: service_action format
permissionSlug := "products_create"

// ❌ Wrong: dot notation
permissionSlug := "products.create"
```

3. **Check Role Preloading:**
```go
// Ensure roles AND permissions are loaded
err = facades.Orm().Query().
    Where("id = ?", user.ID).
    With("Roles.Permissions").  // Critical: Load the relationship!
    First(&userWithRoles)
```

3. **Database Schema Out of Sync:**
   ```bash
   # Run fresh migrations after model changes
   go run . artisan migrate:fresh
   go run . artisan db:seed --seeder=DatabaseSeeder
   ```

4. **Permission Errors After Code Changes:**
   ```bash
   # Re-seed RBAC permissions
   go run . artisan seed --seeder=rbac
   ```

### 🔧 Permission System Issues

**Issue: Write operations (Create/Edit/Delete) not showing in UI**

This is often caused by permission mapping mismatches. Check these common issues:

1. **Resource Name Mismatch:**
```go
// ❌ Wrong - singular resource name
permissions := c.BuildPermissionsMap(ctx, "book")

// ✅ Correct - plural resource name matching RBAC permissions
permissions := c.BuildPermissionsMap(ctx, "books")
```

2. **User Roles Not Loaded:**
```go
// ❌ Wrong - roles not preloaded
err = facades.Orm().Query().
    Where("id = ?", user.ID).
    First(&userWithRoles)

// ✅ Correct - preload roles relationship
err = facades.Orm().Query().
    Where("id = ?", user.ID).
    With("Roles").  // Critical: Load roles!
    First(&userWithRoles)
```

3. **Permission Format Mismatch:**
```go
// ❌ Wrong in controller
if err := c.CheckPermission(ctx, "create.books", nil); err != nil

// ✅ Correct - matches RBAC seeder format
if err := c.CheckPermission(ctx, "books.create", nil); err != nil
```

4. **Database Column Name Issues:**
```go
// Check your user_roles migration has nullable fields:
table.UnsignedBigInteger("assigned_by_id").Nullable()
table.Timestamp("expires_at").Nullable()

// And ensure column mapping in models:
Note string `gorm:"type:text;column:notes"` // Note the 's' in column name
```

5. **Frontend Permission Props:**
```typescript
// Ensure the page controller passes correct permissions
permissions: {
    canCreate: boolean;    // Maps to books.create
    canEdit: boolean;      // Maps to books.update
    canDelete: boolean;    // Maps to books.delete
    canManageLibrary: boolean; // Maps to books.manage
}
```

**Quick Fix Checklist:**
- [ ] Use plural resource names ("books" not "book")
- [ ] Preload user roles with `.With("Roles")`
- [ ] Match permission format: `resource.action`
- [ ] Check nullable columns in migrations
- [ ] Verify frontend receives permission props

---

## 🎓 Learning Path for Junior Developers

### Phase 1: Understanding the Basics
1. **Generated Files**: Examine each generated file
2. **Contracts**: Understand why contracts prevent bugs
3. **RBAC**: Learn how permissions protect resources
4. **Frontend Integration**: See how backend connects to React

### Phase 2: Customization
1. **Add Custom Fields**: Practice modifying models and migrations
2. **Custom Validation**: Create complex validation rules
3. **Business Logic**: Add service layer customizations
4. **UI Enhancements**: Improve the frontend components

### Phase 3: Advanced Features
1. **Relationships**: Implement complex model relationships
2. **File Uploads**: Add file handling capabilities
3. **Real-time Updates**: Integrate WebSocket notifications
4. **Performance**: Add caching and optimization

---

## 📚 Reference

### Generated Files Checklist

- [ ] **Model** with soft deletes and validation
- [ ] **Migration** with indexes and constraints
- [ ] **Service** implementing all contracts
- [ ] **Validation Requests** for create/update
- [ ] **API Controller** with permissions
- [ ] **Page Controller** for Inertia.js
- [ ] **Routes** with middleware
- [ ] **Permissions** for RBAC
- [ ] **TypeScript Types** for frontend
- [ ] **React Components** for UI
- [ ] **React Pages** for full functionality

### Commands Reference

```bash
# Generate complete CRUD
go run . artisan make:crud-e2e ResourceName

# Database operations
go run . artisan migrate
go run . artisan migrate:rollback

# RBAC operations
go run . artisan rbac:setup
go run . artisan rbac:assign email@domain.com role-name

# Seeding
go run . artisan seed --seeder=rbac
```

### Contract Interfaces

All services must implement:
- `CrudServiceContract` - Basic CRUD operations
- `PaginationServiceContract` - Pagination support
- `SortableServiceContract` - Sorting capabilities
- `FilterableServiceContract` - Filtering and search

All controllers must implement:
- `CrudControllerContract` - Standard CRUD endpoints
- `PermissionControllerContract` - Permission enforcement

---

## ✅ Success Criteria

Your CRUD system is properly implemented when:

- [ ] **Backend API** responds to all CRUD operations
- [ ] **Permissions** are enforced on all endpoints
- [ ] **Frontend UI** displays data with proper pagination/sorting
- [ ] **Form Validation** works on both client and server
- [ ] **Error Handling** provides meaningful feedback
- [ ] **Type Safety** prevents runtime errors
- [ ] **Database** properly indexes and constrains data
- [ ] **Soft Deletes** preserve data integrity

Congratulations! You now have a production-ready, secure, and maintainable CRUD system! 🎉