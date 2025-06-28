# Enhanced CRUD System with Repository Pattern & Authorization

This document describes the enhanced CRUD system that combines our simplified approach with repository pattern, validation, and authorization using Goravel's native capabilities.

## Architecture Overview

```
┌─────────────────────┐
│    Controllers      │ ← HTTP layer with validation & authorization 
├─────────────────────┤
│ Validation Requests │ ← Goravel request validation classes
├─────────────────────┤
│     Services        │ ← Business logic layer
├─────────────────────┤
│   Repositories      │ ← Data access abstraction
├─────────────────────┤
│ Auth & Gate Helpers │ ← Authorization logic
├─────────────────────┤
│      Models         │ ← Database entities
└─────────────────────┘
```

## Key Components

### 1. Repository Pattern (`app/repositories/`)

**Base Repository** provides common CRUD operations:
```go
type Repository interface {
    Find(id interface{}) (interface{}, error)
    FindMany(ids []interface{}) ([]interface{}, error)
    Create(data map[string]interface{}) (interface{}, error)
    Update(id interface{}, data map[string]interface{}) (interface{}, error)
    Delete(id interface{}) error
    Count(conditions map[string]interface{}) (int64, error)
    Exists(id interface{}) bool
    BulkCreate(data []map[string]interface{}) ([]interface{}, error)
    BulkUpdate(conditions map[string]interface{}, data map[string]interface{}) error
    BulkDelete(conditions map[string]interface{}) error
}
```

**Queryable Repository** adds advanced querying:
```go
repo := NewQueryableRepository(model, "table_name")
result, err := repo.
    Where("status", "=", "AVAILABLE").
    Search(term, []string{"title", "author"}).
    OrderBy("title", "ASC").
    Paginate(1, 20)
```

**Searchable Repository** adds full-text search:
```go
repo := NewSearchableRepository(model, "table_name")
result, err := repo.
    SearchWithFilters(term, searchFields, filters).
    Paginate(1, 20)
```

### 2. Validation Requests (`app/http/requests/`)

**Create Request** with comprehensive validation:
```go
type BookCreateRequest struct {
    Title       string   `form:"title" json:"title"`
    Author      string   `form:"author" json:"author"`
    ISBN        string   `form:"isbn" json:"isbn"`
    Price       float64  `form:"price" json:"price"`
    // ... other fields
}

func (r *BookCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "title":  "required|max:255",
        "author": "required|max:100",
        "isbn":   "required|regex:^[0-9]{10,13}$|unique:books,isbn",
        "price":  "required|numeric|min:0",
    }
}
```

**Update Request** with partial validation:
```go
type BookUpdateRequest struct {
    Title  *string  `form:"title" json:"title"`   // Pointers for optional fields
    Author *string  `form:"author" json:"author"`
    Price  *float64 `form:"price" json:"price"`
    // ... other fields
}
```

### 3. Authorization Gates (`app/providers/gate_service_provider.go`)

**Resource-based permissions**:
```go
// Standard CRUD gates
facades.Gate().Define("viewAny.books", func(ctx http.Context, user interface{}) access.Response {
    // Everyone can view book lists
    return access.NewAllowResponse()
})

facades.Gate().Define("create.books", func(ctx http.Context, user interface{}) access.Response {
    // Only moderators and admins can create
    return gateHelper.RoleBasedAccess("ADMIN", "MODERATOR")(ctx, user)
})

facades.Gate().Define("delete.books", func(ctx http.Context, user interface{}, model interface{}) access.Response {
    // Only admins can delete
    return gateHelper.RoleBasedAccess("ADMIN")(ctx, user)
})
```

**Domain-specific gates**:
```go
facades.Gate().Define("borrow.books", func(ctx http.Context, user interface{}) access.Response {
    // All authenticated users can borrow books
    return authHelper.IsAuthenticated(ctx) ? access.NewAllowResponse() : access.NewDenyResponse("Must be authenticated")
})
```

### 4. Enhanced Services (`app/services/`)

**Repository Integration**:
```go
type BookService struct {
    repository contracts.SearchableRepository
    authHelper contracts.AuthHelper
}

func (s *BookService) GetList(req helpers.ListRequest) (*contracts.PaginatedResult, error) {
    repo := s.repository.Reset()
    
    if req.Search != "" {
        repo = repo.Search(req.Search, book.SearchFields())
    }
    
    return repo.OrderBy("id", "DESC").Paginate(req.Page, req.PageSize)
}
```

### 5. Enhanced Controllers (`app/http/controllers/`)

**Authorization Checks**:
```go
func (c *BookController) Store(ctx http.Context) http.Response {
    // Check authorization
    if response := facades.Gate().Inspect("create.books", ctx); response.Denied() {
        return ctx.Response().Json(http.StatusForbidden, map[string]string{
            "error": response.Message(),
        })
    }
    
    // Validate request
    var createRequest requests.BookCreateRequest
    if errors, err := ctx.Request().ValidateRequest(&createRequest); err != nil || errors != nil {
        // Handle validation errors
    }
    
    // Create using validated data
    book, err := c.bookService.Create(createRequest.ToCreateData())
    // ...
}
```

## Usage Examples

### Creating a New Resource

#### 1. Define the Model
```go
type Product struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Price       float64   `json:"price" gorm:"not null"`
    CategoryID  uint      `json:"categoryId"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

func (p Product) SearchFields() []string {
    return []string{"name", "description"}
}
```

#### 2. Create Repository
```go
type ProductRepository struct {
    *repositories.SearchableBaseRepository
}

func NewProductRepository() contracts.SearchableRepository {
    return &ProductRepository{
        SearchableBaseRepository: repositories.NewSearchableRepository(&models.Product{}, "products"),
    }
}
```

#### 3. Create Validation Requests
```go
type ProductCreateRequest struct {
    Name        string  `form:"name" json:"name"`
    Price       float64 `form:"price" json:"price"`
    CategoryID  uint    `form:"categoryId" json:"categoryId"`
}

func (r *ProductCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "name":       "required|max:255",
        "price":      "required|numeric|min:0",
        "categoryId": "required|exists:categories,id",
    }
}
```

#### 4. Create Service
```go
type ProductService struct {
    repository contracts.SearchableRepository
    authHelper contracts.AuthHelper
}

func NewProductService() *ProductService {
    return &ProductService{
        repository: NewProductRepository(),
        authHelper: helpers.NewAuthHelper(),
    }
}
```

#### 5. Create Controller
```go
type ProductController struct {
    productService *ProductService
}

func (c *ProductController) Store(ctx http.Context) http.Response {
    if response := facades.Gate().Inspect("create.products", ctx); response.Denied() {
        return ctx.Response().Json(http.StatusForbidden, map[string]string{
            "error": response.Message(),
        })
    }
    
    var createRequest ProductCreateRequest
    if errors, err := ctx.Request().ValidateRequest(&createRequest); err != nil || errors != nil {
        // Handle validation
    }
    
    product, err := c.productService.Create(createRequest.ToCreateData())
    return ctx.Response().Json(http.StatusCreated, product)
}
```

#### 6. Register Gates
```go
func RegisterProductGates() {
    gateHelper := helpers.NewGateHelper()
    
    config := contracts.GateConfig{
        CreateHandler: gateHelper.RoleBasedAccess("ADMIN", "MANAGER"),
        UpdateHandler: gateHelper.RoleBasedAccess("ADMIN", "MANAGER"),
        DeleteHandler: gateHelper.RoleBasedAccess("ADMIN"),
    }
    
    gateHelper.RegisterResourceGates("products", config)
}
```

## API Usage Examples

### Advanced Querying
```bash
# Basic list with search
GET /books?search=golang&page=1&pageSize=20&sort=title ASC

# Advanced filtering
GET /books/advanced?status=AVAILABLE&author=Martin&minPrice=20&maxPrice=100

# Pagination with repository
GET /books?page=2&pageSize=50
```

### Validation Examples
```bash
# Create with validation
POST /books
{
    "title": "Clean Code",
    "author": "Robert Martin",
    "isbn": "9780132350884",
    "price": 45.99
}

# Validation error response
{
    "error": "Validation failed",
    "errors": {
        "isbn": ["This ISBN already exists"],
        "price": ["Price must be greater than or equal to 0"]
    }
}
```

### Authorization Examples
```bash
# Unauthorized access
GET /books
Response: 403 Forbidden
{
    "error": "Insufficient role privileges"
}

# Successful authorization
POST /books/123/borrow
Response: 200 OK
{
    "message": "Book borrowed successfully"
}
```

## Key Benefits

### 1. **Repository Pattern Benefits**
- **Data Access Abstraction**: Clean separation between business logic and data access
- **Testability**: Easy to mock repositories for unit testing
- **Flexibility**: Can switch between different data sources
- **Query Building**: Fluent interface for complex queries

### 2. **Validation Benefits**
- **Goravel Native**: Uses Goravel's built-in validation system
- **Reusable**: Request classes can be reused across controllers
- **Comprehensive**: Supports database validation (unique, exists)
- **Custom Messages**: Localized and custom error messages

### 3. **Authorization Benefits**
- **Gate System**: Leverages Goravel's Gate system
- **Resource-based**: Permissions tied to specific resources
- **Role-based**: Simple role hierarchy support
- **Flexible**: Custom authorization logic when needed

### 4. **Maintained Simplicity**
- **No Over-engineering**: Still practical and straightforward
- **Progressive Enhancement**: Add features as needed
- **Go Idiomatic**: Follows Go conventions
- **Clear Architecture**: Easy to understand and maintain

## Migration from Simple CRUD

To migrate existing simple CRUD to enhanced version:

1. **Create Repository**: Replace direct database calls with repository pattern
2. **Add Validation**: Create request classes for validation
3. **Add Authorization**: Define gates for resource permissions
4. **Update Service**: Inject repository and auth helper
5. **Update Controller**: Add gate checks and request validation

This enhanced system provides enterprise-level features while maintaining the simplicity and practicality of our original approach.