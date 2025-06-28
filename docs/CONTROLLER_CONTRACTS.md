# Controller Contracts Implementation

This document describes the comprehensive controller contract system that enforces CRUD operations, pagination, validation, and proper response handling.

## Overview

The controller contract system ensures that:
- **CRUD Operations**: All controllers MUST implement Create, Read, Update, Delete
- **Pagination**: All listing endpoints MUST implement proper pagination
- **Validation**: All data input MUST be validated using contracts
- **Response Format**: All responses MUST follow standardized format
- **Authorization**: All endpoints MUST implement permission checking

## Contracts Implemented

### 1. CrudControllerContract
Forces implementation of core CRUD operations:
```go
type CrudControllerContract interface {
    Index(ctx http.Context) http.Response  // GET /resource - List with pagination
    Show(ctx http.Context) http.Response   // GET /resource/{id} - JSON for modals
    Store(ctx http.Context) http.Response  // POST /resource - Create
    Update(ctx http.Context) http.Response // PUT /resource/{id} - Update
    Delete(ctx http.Context) http.Response // DELETE /resource/{id} - Delete
    
    PaginationControllerContract
    ValidationControllerContract
    ResponseControllerContract
}
```

### 2. PaginationControllerContract
Enforces pagination in all listing endpoints:
```go
type PaginationControllerContract interface {
    ValidatePaginationRequest(ctx http.Context) (*ListRequest, error)
    GetPaginationDefaults() (page int, pageSize int, maxPageSize int)
    BuildPaginatedResponse(result *PaginatedResult, request *ListRequest) map[string]interface{}
}
```

### 3. ValidationControllerContract
Enforces validation for Create/Update operations:
```go
type ValidationControllerContract interface {
    ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error)
    ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error)
    ValidateID(ctx http.Context, paramName string) (uint, error)
    GetValidationRules() map[string]interface{}
}
```

### 4. ResponseControllerContract
Enforces consistent response formatting:
```go
type ResponseControllerContract interface {
    SuccessResponse(ctx http.Context, data interface{}, message string) http.Response
    CreatedResponse(ctx http.Context, data interface{}, message string) http.Response
    ValidationErrorResponse(ctx http.Context, errors map[string]interface{}) http.Response
    ResourceNotFoundResponse(ctx http.Context, resourceType string, id uint) http.Response
    // ... more response methods
}
```

### 5. AuthorizationControllerContract
Enforces authorization checks:
```go
type AuthorizationControllerContract interface {
    CheckPermission(ctx http.Context, permission string, resource interface{}) error
    GetCurrentUser(ctx http.Context) interface{}
    RequireAuthentication(ctx http.Context) error
    BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool
}
```

### 6. PageControllerContract
For Inertia.js page controllers:
```go
type PageControllerContract interface {
    Index(ctx http.Context) http.Response
    PaginationControllerContract
    AuthorizationControllerContract
    PageResponseContract
}
```

## Implementation Examples

### Resource Controller (API)
```go
type BookController struct {
    *contracts.BaseCrudController
    bookService *services.BookService
    authHelper  contracts.AuthHelper
}

func NewBookController() *BookController {
    controller := &BookController{
        BaseCrudController: contracts.NewBaseCrudController("book"),
        bookService:        services.NewBookService(),
        authHelper:         helpers.NewAuthHelper(),
    }

    // Automatic validation - panics if contracts not implemented
    contracts.MustRegisterCrudController("books", controller)
    return controller
}

// All CRUD methods are enforced by interface
func (c *BookController) Index(ctx http.Context) http.Response {
    // Contract-enforced pagination validation
    req, err := c.ValidatePaginationRequest(ctx)
    if err != nil {
        return c.BadRequestResponse(ctx, "Invalid pagination", nil)
    }

    result, err := c.bookService.GetList(*req)
    if err != nil {
        return c.InternalErrorResponse(ctx, err.Error())
    }

    // Contract-enforced response format
    response := c.BuildPaginatedResponse(result, req)
    return c.SuccessResponse(ctx, response, "Books retrieved successfully")
}

func (c *BookController) Show(ctx http.Context) http.Response {
    // Contract-enforced ID validation
    id, err := c.ValidateID(ctx, "id")
    if err != nil {
        return c.BadRequestResponse(ctx, "Invalid ID", nil)
    }

    // Contract-enforced authorization
    if err := c.CheckPermission(ctx, "view.books", nil); err != nil {
        return c.ForbiddenResponse(ctx, "Access denied")
    }

    book, err := c.bookService.GetByID(id)
    if err != nil {
        return c.ResourceNotFoundResponse(ctx, "book", id)
    }

    // Returns JSON for modal display
    return c.SuccessResponse(ctx, book, "Book details retrieved")
}
```

### Page Controller (Inertia.js)
```go
type BooksPageController struct {
    *contracts.BasePageController
    bookService *services.BookService
    authHelper  contracts.AuthHelper
}

func NewBooksPageController() *BooksPageController {
    controller := &BooksPageController{
        BasePageController: contracts.NewBasePageController("book", "Books/Index"),
        bookService:        services.NewBookService(),
        authHelper:         helpers.NewAuthHelper(),
    }

    // Automatic validation - panics if contracts not implemented
    contracts.MustRegisterPageController("books_page", controller)
    return controller
}

func (c *BooksPageController) Index(ctx http.Context) http.Response {
    // Contract-enforced page request validation
    req, err := c.ValidatePageRequest(ctx)
    if err != nil {
        req = &contracts.ListRequest{Page: 1, PageSize: 20}
        req.SetDefaults()
    }

    // Contract-enforced permissions map
    permissions := c.BuildPermissionsMap(ctx, "book")

    result, err := c.bookService.GetList(*req)
    if err != nil {
        // Graceful error handling for pages
        result = &contracts.PaginatedResult{Data: []interface{}{}}
    }

    // Contract-enforced page props structure
    props := c.BuildPageProps(result.Data, req, permissions, nil)
    return inertia.Render(ctx, "Books/Index", props)
}
```

## Key Features

### 1. **Impossible to Skip Requirements**
Controllers CANNOT be instantiated without implementing all required methods:
```go
// This will PANIC at startup if any method is missing
contracts.MustRegisterCrudController("books", controller)
```

### 2. **Enforced Pagination**
All listing endpoints MUST handle pagination properly:
```go
// Validates page, pageSize, sort, direction automatically
req, err := c.ValidatePaginationRequest(ctx)
```

### 3. **Standardized Responses**
All responses follow consistent format:
```go
// Success responses
c.SuccessResponse(ctx, data, "Operation successful")
c.CreatedResponse(ctx, resource, "Resource created")

// Error responses  
c.ValidationErrorResponse(ctx, errors)
c.ResourceNotFoundResponse(ctx, "book", id)
```

### 4. **Automatic Validation**
Request validation is enforced by contracts:
```go
// Validates and returns clean data or error
data, err := c.ValidateCreateRequest(ctx)
data, err := c.ValidateUpdateRequest(ctx, id)
```

### 5. **JSON for Modals**
Show endpoints return JSON specifically for modal display, not full pages.

## Response Formats

### Standard API Response
```json
{
  "success": true,
  "data": { /* resource data */ },
  "message": "Operation successful",
  "meta": {
    "version": "1.0.0",
    "timestamp": 1234567890
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [ /* array of resources */ ],
  "pagination": {
    "current_page": 1,
    "last_page": 5,
    "per_page": 20,
    "total": 100,
    "has_next": true,
    "has_prev": false
  },
  "filters": {
    "page": 1,
    "pageSize": 20,
    "search": "query",
    "sort": "created_at",
    "direction": "DESC"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "title": ["Title is required"],
    "email": ["Email format is invalid"]
  }
}
```

## Benefits

1. **Contract Enforcement**: Impossible to create incomplete controllers
2. **Consistent API**: All endpoints follow same patterns
3. **Automatic Validation**: Built-in request/response validation
4. **Pagination Enforcement**: No endpoint can skip pagination
5. **Error Handling**: Standardized error responses
6. **Authorization**: Enforced permission checking
7. **Modal Support**: Show endpoints return JSON for modal display
8. **Type Safety**: All operations use proper types
9. **Extensible**: Easy to add new contracts or requirements

The system ensures that any developer creating a new controller **cannot skip** implementing proper CRUD operations, pagination, validation, or response handling - it's enforced by the interface contracts and factory validation!