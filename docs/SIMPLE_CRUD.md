# Simple Goravel CRUD System

A lightweight, practical approach to CRUD operations with built-in pagination, sorting, and filtering - without heavy DDD concepts.

## Core Philosophy

- **Simple over Complex**: Favor straightforward code over abstract patterns
- **Convention over Configuration**: Use sensible defaults
- **Minimal Boilerplate**: Maximum functionality with minimal code
- **Go Idiomatic**: Follow Go conventions, not enterprise patterns

## Architecture

```
┌─────────────────┐
│   Controllers   │ ← Simple HTTP handlers with direct service calls
├─────────────────┤
│    Services     │ ← Business logic with CrudHelper for database ops
├─────────────────┤
│   CrudHelper    │ ← Simple database operations wrapper
├─────────────────┤
│  FilterBuilder  │ ← Optional advanced query building
├─────────────────┤
│     Models      │ ← Plain structs with simple methods
└─────────────────┘
```

## Core Components

### 1. CrudHelper (`app/helpers/crud_helper.go`)

Simple wrapper for database operations:

```go
type CrudHelper struct {
    tableName string
    db        orm.Query
}

func NewCrudHelper(tableName string) *CrudHelper {
    return &CrudHelper{
        tableName: tableName,
        db:        facades.Orm().Query().Table(tableName),
    }
}
```

### 2. Request/Response Types (`app/helpers/types.go`)

Simple, practical types:

```go
type ListRequest struct {
    Page     int    `form:"page" json:"page"`
    PageSize int    `form:"pageSize" json:"pageSize"`
    Sort     string `form:"sort" json:"sort"`
    Search   string `form:"search" json:"search"`
}

type ListResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"pageSize"`
    TotalPages int         `json:"totalPages"`
}
```

### 3. FilterBuilder (`app/helpers/filter_builder.go`)

Optional advanced query building:

```go
builder := NewFilterBuilder("books")
    .Search([]string{"title", "author"}, "golang")
    .Where("status", "AVAILABLE")
    .WhereRange("price", 10.0, 50.0)
    .Order("title ASC")
    .Paginate(1, 20)

result := builder.Build()
```

## Book Example Implementation

### Model (`app/models/book.go`)

Just a regular struct - no inheritance or complex interfaces:

```go
type Book struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Author      string    `json:"author" gorm:"not null"`
    ISBN        string    `json:"isbn" gorm:"unique;not null"`
    Description string    `json:"description"`
    Price       float64   `json:"price" gorm:"default:0"`
    Status      string    `json:"status" gorm:"default:'AVAILABLE'"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

// Optional: Define searchable fields
func (b Book) SearchFields() []string {
    return []string{"title", "author", "isbn", "description"}
}
```

### Service (`app/services/book_service.go`)

Simple, focused business logic:

```go
type BookService struct {
    crud *helpers.CrudHelper
}

func NewBookService() *BookService {
    return &BookService{
        crud: helpers.NewCrudHelper("books"),
    }
}

// Simple list with pagination and search
func (s *BookService) GetList(req helpers.ListRequest) (*helpers.ListResponse, error) {
    req.SetDefaults()
    
    book := models.Book{}
    builder := helpers.NewFilterBuilder("books")
    
    if req.Search != "" {
        builder = builder.Search(book.SearchFields(), req.Search)
    }
    
    total, err := builder.Count()
    if err != nil {
        return nil, err
    }
    
    query := builder.Order(req.Sort).Paginate(req.Page, req.PageSize).Build()
    
    var books []models.Book
    err = query.Find(&books)
    if err != nil {
        return nil, err
    }
    
    totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
    
    return &helpers.ListResponse{
        Data:       books,
        Total:      total,
        Page:       req.Page,
        PageSize:   req.PageSize,
        TotalPages: totalPages,
    }, nil
}

// Simple CRUD operations
func (s *BookService) GetByID(id uint) (*models.Book, error) {
    var book models.Book
    err := s.crud.Query().Where("id", id).First(&book)
    return &book, err
}

func (s *BookService) Create(data map[string]interface{}) (*models.Book, error) {
    if err := s.validateBookData(data, false); err != nil {
        return nil, err
    }
    
    err := facades.Orm().Query().Create(&data)
    if err != nil {
        return nil, err
    }
    
    // Return the created book
    if id, ok := data["id"].(uint); ok {
        return s.GetByID(id)
    }
    return nil, nil
}
```

### Controller (`app/http/controllers/book_controller.go`)

Clean, simple HTTP handlers:

```go
type BookController struct {
    bookService *services.BookService
}

func NewBookController() *BookController {
    return &BookController{
        bookService: services.NewBookService(),
    }
}

// GET /books
func (c *BookController) Index(ctx http.Context) http.Response {
    var req helpers.ListRequest
    ctx.Request().Bind(&req)
    
    result, err := c.bookService.GetList(req)
    if err != nil {
        return ctx.Response().Json(500, map[string]string{"error": err.Error()})
    }
    
    return ctx.Response().Json(200, result)
}

// POST /books
func (c *BookController) Store(ctx http.Context) http.Response {
    var data map[string]interface{}
    ctx.Request().Bind(&data)
    
    book, err := c.bookService.Create(data)
    if err != nil {
        return ctx.Response().Json(422, map[string]string{"error": err.Error()})
    }
    
    return ctx.Response().Json(201, book)
}
```

## API Usage Examples

### Basic Operations

```bash
# List books with pagination and search
GET /books?page=1&pageSize=20&sort=title ASC&search=golang

# Advanced filtering
GET /books/advanced?status=AVAILABLE&author=Martin&minPrice=20&maxPrice=100

# Get single book
GET /books/123

# Create book
POST /books
{
    "title": "Clean Code",
    "author": "Robert Martin", 
    "isbn": "9780132350884",
    "price": 45.99
}

# Update book
PUT /books/123
{
    "price": 39.99,
    "status": "BORROWED"
}

# Delete book (soft delete)
DELETE /books/123

# Domain-specific operations
POST /books/123/borrow
POST /books/123/return
```

### Response Format

```json
{
    "data": [
        {
            "id": 1,
            "title": "Clean Code",
            "author": "Robert Martin",
            "isbn": "9780132350884", 
            "price": 45.99,
            "status": "AVAILABLE",
            "createdAt": "2024-01-01T00:00:00Z"
        }
    ],
    "total": 150,
    "page": 1,
    "pageSize": 20,
    "totalPages": 8
}
```

## Key Benefits

### 1. **Dramatically Less Code**
- No complex interfaces or inheritance hierarchies
- Direct, readable implementations
- 80% less boilerplate than enterprise patterns

### 2. **Easy to Understand**
- Simple struct compositions
- No hidden abstractions
- Clear data flow

### 3. **Still Powerful**
- Pagination, sorting, filtering out of the box
- Advanced query building when needed
- Domain-specific operations

### 4. **Go Idiomatic**
- Follows Go conventions
- Explicit error handling
- Simple composition over inheritance

### 5. **Maintainable**
- Easy to debug and modify
- No complex dependency injection
- Straightforward testing

## Creating New Resources

To add a new resource (e.g., Product), simply:

1. **Create the model** with optional `SearchFields()` method
2. **Create the service** using `CrudHelper`
3. **Create the controller** with simple HTTP handlers
4. **Register routes** in a new routes file

Example for Product:

```go
// app/services/product_service.go
type ProductService struct {
    crud *helpers.CrudHelper
}

func NewProductService() *ProductService {
    return &ProductService{
        crud: helpers.NewCrudHelper("products"),
    }
}

// Use the same patterns as BookService...
```

## Migration

```bash
# Run migrations to create the books table
go run . artisan migrate

# Create sample data
go run . artisan db:seed
```

## Testing

Simple testing with no mocks needed:

```go
func TestBookService_Create(t *testing.T) {
    service := services.NewBookService()
    
    data := map[string]interface{}{
        "title": "Test Book",
        "author": "Test Author",
        "isbn": "1234567890",
    }
    
    book, err := service.Create(data)
    assert.NoError(t, err)
    assert.Equal(t, "Test Book", book.Title)
}
```

This simplified approach provides 90% of the benefits with 50% of the complexity. Perfect for most real-world applications that need practical CRUD operations without enterprise overhead.