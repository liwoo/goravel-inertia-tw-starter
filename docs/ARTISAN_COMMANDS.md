# Custom Artisan Commands for CRUD System

This document describes the custom Artisan commands for generating CRUD resources with our enhanced system.

## Available Commands

### 1. `make:service` - Create Service Class

Creates a new service class with repository pattern integration.

```bash
go run . artisan make:service {name}
```

**Example:**
```bash
go run . artisan make:service PlayerService
go run . artisan make:service Product  # Automatically adds "Service" suffix
```

**Generates:**
- `app/services/player_service.go` - Service with repository integration
- Includes: GetList, GetByID, Create, Update, Delete methods
- Built-in validation and error handling
- Repository pattern integration

---

### 2. `make:request` - Create Validation Request

Creates validation request classes for form validation.

```bash
go run . artisan make:request {name} [--update]
```

**Examples:**
```bash
# Create request (for creation)
go run . artisan make:request PlayerCreateRequest

# Create update request (with pointer fields for partial updates)
go run . artisan make:request Player --update
```

**Generates:**
- `app/http/requests/player_create_request.go` - For creation validation
- `app/http/requests/player_update_request.go` - For update validation (with --update flag)
- Includes: Rules, Messages, Attributes, Authorization methods
- Data transformation methods (ToCreateData, ToUpdateData)

---

### 3. `make:repository` - Create Repository Class

Creates a repository class with searchable capabilities.

```bash
go run . artisan make:repository {name}
```

**Example:**
```bash
go run . artisan make:repository PlayerRepository
go run . artisan make:repository Product  # Automatically adds "Repository" suffix
```

**Generates:**
- `app/repositories/player_repository.go` - Repository with search capabilities
- Extends SearchableBaseRepository
- Includes domain-specific methods (GetByName, GetByStatus, etc.)
- Built-in search and filtering

---

### 4. `make:crud` - Complete CRUD Generator

Creates a complete CRUD resource with all components.

```bash
go run . artisan make:crud {name} [--model] [--migration] [--routes]
```

**Options:**
- `--model` - Also create model file
- `--migration` - Also create migration file  
- `--routes` - Show route examples after generation

**Example:**
```bash
# Complete CRUD with model and migration
go run . artisan make:crud Player --model --migration --routes

# CRUD without model/migration (if they already exist)
go run . artisan make:crud Team --routes
```

**Generates:**
- `app/models/player.go` (if --model)
- `database/migrations/xxx_create_players_table.go` (if --migration)
- `app/repositories/player_repository.go`
- `app/services/player_service.go`
- `app/http/requests/player_create_request.go`
- `app/http/requests/player_update_request.go`
- `app/http/controllers/player_controller.go`
- `app/providers/player_gate_provider.go` (authorization gates)

## Usage Examples

### Example 1: Creating a Product CRUD

```bash
# Generate complete Product CRUD
go run . artisan make:crud Product --model --migration --routes
```

This creates:
```
app/
├── models/product.go
├── repositories/product_repository.go
├── services/product_service.go
├── http/
│   ├── controllers/product_controller.go
│   └── requests/
│       ├── product_create_request.go
│       └── product_update_request.go
└── providers/product_gate_provider.go

database/
└── migrations/20240101000000_create_products_table.go
```

### Example 2: Creating Individual Components

```bash
# Create just a service
go run . artisan make:service OrderService

# Create just validation requests
go run . artisan make:request OrderCreateRequest
go run . artisan make:request Order --update

# Create just a repository
go run . artisan make:repository OrderRepository
```

## Generated Code Features

### Service Features
- **Repository Integration**: Uses repository pattern for data access
- **Validation**: Built-in validation methods
- **Error Handling**: Comprehensive error handling with proper messages
- **Search & Filtering**: Advanced search and filtering capabilities
- **Pagination**: Built-in pagination support

### Request Features
- **Goravel Validation**: Uses Goravel's native validation system
- **Custom Messages**: Localized and custom error messages
- **Authorization**: Built-in authorization checks
- **Data Transformation**: Clean conversion to service data
- **Partial Updates**: Update requests use pointer fields for optional fields

### Controller Features
- **Authorization Gates**: Integrated with Goravel's Gate system
- **Validation**: Request validation with error handling
- **Consistent Responses**: Standardized JSON responses
- **Error Handling**: Proper HTTP status codes and error messages
- **Advanced Filtering**: Query parameter parsing for filters

### Repository Features
- **Searchable**: Full-text search across multiple fields
- **Queryable**: Advanced query building with method chaining
- **Domain Methods**: Domain-specific methods (GetByName, GetByStatus, etc.)
- **Bulk Operations**: Support for bulk create, update, delete

## Post-Generation Steps

After generating CRUD resources:

1. **Update Model Fields** (if generated):
   ```go
   // app/models/product.go - Add your actual fields
   type Product struct {
       ID          uint      `json:"id" gorm:"primaryKey"`
       Name        string    `json:"name" gorm:"not null"`
       Price       float64   `json:"price" gorm:"not null"`
       CategoryID  uint      `json:"categoryId"`
       // Add more fields...
   }
   ```

2. **Update Validation Rules**:
   ```go
   // app/http/requests/product_create_request.go
   func (r *ProductCreateRequest) Rules(ctx http.Context) map[string]string {
       return map[string]string{
           "name":       "required|max:255",
           "price":      "required|numeric|min:0",
           "categoryId": "required|exists:categories,id",
           // Add more rules...
       }
   }
   ```

3. **Add Repository Methods**:
   ```go
   // app/repositories/product_repository.go
   func (r *ProductRepository) GetByCategory(categoryID uint) ([]models.Product, error) {
       // Add domain-specific methods
   }
   ```

4. **Register Routes**:
   ```go
   // routes/api.go or routes/web.go
   func ProductRoutes() {
       productController := controllers.NewProductController()
       
       facades.Route().Group(func(router route.Router) {
           router.Get("/products", productController.Index)
           router.Get("/products/{id}", productController.Show)
           // ... more routes
       })
   }
   ```

5. **Run Migration**:
   ```bash
   go run . artisan migrate
   ```

6. **Register Gate Provider** (if needed):
   ```go
   // bootstrap/app.go or your service registration
   app.Register(providers.ProductGateServiceProvider{})
   ```

## API Endpoints Generated

For a `Product` resource, the following endpoints are available:

```
GET    /products              - List products with pagination/search
GET    /products/{id}         - Get single product
GET    /products/advanced     - Advanced filtering
POST   /products              - Create product (auth required)
PUT    /products/{id}         - Update product (auth required)  
DELETE /products/{id}         - Delete product (auth required)
```

**Query Parameters:**
- `page` - Page number for pagination
- `pageSize` - Items per page
- `search` - Search term
- `sort` - Sorting (e.g., "name ASC")
- `status` - Filter by status
- Custom filters based on your model

**Example Requests:**
```bash
# List with search and pagination
GET /products?search=laptop&page=1&pageSize=20&sort=name ASC

# Advanced filtering
GET /products/advanced?status=ACTIVE&minPrice=100&maxPrice=500

# Create product
POST /products
{
    "name": "Gaming Laptop",
    "price": 999.99,
    "categoryId": 1
}
```

## Customization

### Adding Custom Service Methods

```go
// app/services/product_service.go
func (s *ProductService) GetPopular(limit int) ([]models.Product, error) {
    // Add domain-specific business logic
    return s.repository.GetPopular(limit)
}
```

### Adding Custom Repository Methods

```go
// app/repositories/product_repository.go
func (r *ProductRepository) GetPopular(limit int) ([]models.Product, error) {
    var products []models.Product
    err := r.builder.
        Where("status", "=", "ACTIVE").
        OrderBy("views", "DESC").
        Limit(limit).
        Find(&products)
    return products, err
}
```

### Adding Custom Controller Endpoints

```go
// app/http/controllers/product_controller.go
func (c *ProductController) GetPopular(ctx http.Context) http.Response {
    limit := c.getIntParam(ctx, "limit", 10)
    products, err := c.productService.GetPopular(limit)
    if err != nil {
        return ctx.Response().Json(500, map[string]string{"error": err.Error()})
    }
    return ctx.Response().Json(200, products)
}
```

This Artisan command system provides a solid foundation for rapid CRUD development while maintaining code quality and consistency across your application.