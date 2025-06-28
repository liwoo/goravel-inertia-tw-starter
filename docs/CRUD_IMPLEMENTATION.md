# CRUD Implementation Guide

This document provides a comprehensive guide to the CRUD system implemented in this Goravel application.

## Overview

The CRUD system follows the contracts defined in `app/docs/CRUD_SVC_CONTRACTS.md` and provides a layered architecture for building scalable, maintainable APIs.

## Architecture

```
┌─────────────────┐
│   Controllers   │ ← HTTP layer, request/response handling
├─────────────────┤
│    Services     │ ← Business logic, validation, operations
├─────────────────┤
│  Repositories   │ ← Data access layer (optional)
├─────────────────┤
│     Models      │ ← Database entities with interfaces
└─────────────────┘
```

## Components

### 1. Contracts (`app/contracts/`)

Define interfaces for consistent behavior across the application:

- **crud.go**: Core CRUD service and controller interfaces
- **repository.go**: Data access layer interfaces  
- **entity.go**: Entity behavior contracts
- **validation.go**: Validation and authorization interfaces
- **config.go**: Configuration and factory interfaces

### 2. Services (`app/services/`)

Business logic layer with progressive enhancement:

- **BaseCrudService**: Core CRUD operations
- **FilterableService**: Adds advanced filtering capabilities
- **SoftDeleteService**: Adds soft delete functionality
- **BookService**: Domain-specific example implementation

### 3. Controllers (`app/http/controllers/base/`)

HTTP layer with consistent response handling:

- **CrudController**: Basic CRUD endpoints
- **SoftDeleteController**: Adds soft delete endpoints
- **BookController**: Domain-specific example

### 4. Utilities (`app/utils/`)

- **ResponseBuilder**: Consistent API response formatting

## Usage Examples

### Creating a New Resource

#### 1. Define the Model

```go
// app/models/product.go
package models

import (
    "github.com/goravel/framework/database/orm"
    "players/app/contracts"
)

type Product struct {
    orm.Model
    Name        string  `json:"name" gorm:"not null"`
    Description string  `json:"description"`
    Price       float64 `json:"price" gorm:"not null"`
    Stock       int     `json:"stock" gorm:"default:0"`
    CategoryID  uint    `json:"category_id"`
    orm.SoftDeletes
}

// Implement Entity interfaces
func (p *Product) GetID() interface{} { return p.ID }
func (p *Product) GetTableName() string { return "products" }
func (p *Product) GetSearchableFields() []string {
    return []string{"name", "description"}
}
func (p *Product) GetFilterableFields() map[string]string {
    return map[string]string{
        "category": "category_id",
        "price_min": "price",
        "price_max": "price",
        "in_stock": "stock",
    }
}
```

#### 2. Create the Service

```go
// app/services/product_service.go
package services

import "players/app/contracts"

type ProductService struct {
    *SoftDeleteService
}

func NewProductService() *ProductService {
    config := &contracts.CrudConfig{
        TableName:       "products",
        PrimaryKey:      "id", 
        DefaultPageSize: 20,
        MaxPageSize:     100,
        DefaultSort:     "created_at DESC",
        SearchFields:    []string{"name", "description"},
        FilterableFields: map[string]string{
            "category": "category_id",
            "price_min": "price",
            "price_max": "price", 
        },
        SoftDelete: true,
        Timestamps: true,
    }

    return &ProductService{
        SoftDeleteService: NewSoftDeleteService(&models.Product{}, config),
    }
}

// Add domain-specific methods
func (s *ProductService) GetLowStock(threshold int) (*contracts.ListResponse, error) {
    filters := map[string]interface{}{
        "stock": map[string]interface{}{"max": threshold},
    }
    
    req := contracts.ListRequest{Page: 1, PageSize: 50}
    return s.GetListWithFilters(req, filters)
}
```

#### 3. Create the Controller

```go
// app/http/controllers/product_controller.go
package controllers

import (
    "github.com/goravel/framework/contracts/http"
    "players/app/http/controllers/base"
    "players/app/services"
)

type ProductController struct {
    *base.SoftDeleteController
    productService *services.ProductService
}

func NewProductController() *ProductController {
    service := services.NewProductService()
    return &ProductController{
        SoftDeleteController: base.NewSoftDeleteController(service, "product"),
        productService:      service,
    }
}

// Add custom endpoints
func (c *ProductController) GetLowStock(ctx http.Context) http.Response {
    threshold := c.getIntParam(ctx, "threshold", 10)
    
    result, err := c.productService.GetLowStock(threshold)
    if err != nil {
        return c.responseBuilder.Error("Failed to get low stock products", 500)
    }
    
    return c.responseBuilder.Paginated(result)
}
```

#### 4. Register Routes

```go
// routes/products.go
package routes

import (
    "github.com/goravel/framework/contracts/route"
    "github.com/goravel/framework/facades"
    "players/app/http/controllers"
)

func ProductRoutes() {
    controller := controllers.NewProductController()
    
    facades.Route().Group(func(router route.Router) {
        // Standard CRUD
        router.Get("/products", controller.Index)
        router.Get("/products/{id}", controller.Show)
        router.Post("/products", controller.Store)
        router.Put("/products/{id}", controller.Update)
        router.Delete("/products/{id}", controller.Destroy)
        
        // Soft delete
        router.Get("/products/trash", controller.Trash)
        router.Post("/products/{id}/restore", controller.Restore)
        router.Delete("/products/{id}/force", controller.ForceDestroy)
        
        // Custom endpoints
        router.Get("/products/low-stock", controller.GetLowStock)
    })
}
```

## API Usage Examples

### Basic CRUD Operations

```bash
# List products with pagination
GET /products?page=1&pageSize=20&sort=name ASC

# Search products
GET /products?search=laptop

# Filter products
GET /products?filter_category=electronics&filter_price_min=100&filter_price_max=500

# Get single product
GET /products/123

# Create product
POST /products
{
    "name": "Laptop", 
    "description": "Gaming laptop",
    "price": 999.99,
    "stock": 50,
    "category_id": 1
}

# Update product
PUT /products/123
{
    "price": 899.99,
    "stock": 45
}

# Soft delete
DELETE /products/123

# List trashed items
GET /products/trash

# Restore
POST /products/123/restore

# Permanent delete
DELETE /products/123/force
```

### Response Format

All responses follow a consistent format:

```json
{
    "success": true,
    "data": {
        "id": 123,
        "name": "Laptop",
        "price": 999.99
    },
    "message": "Resource created successfully"
}
```

Paginated responses include pagination metadata:

```json
{
    "success": true,
    "data": [...],
    "pagination": {
        "total": 150,
        "page": 1,
        "pageSize": 20,
        "totalPages": 8,
        "hasNext": true,
        "hasPrev": false
    }
}
```

## Advanced Features

### Custom Filtering

The system supports complex filtering operations:

```bash
# Range filters
GET /products?filter_price_min=100&filter_price_max=500

# Multiple values (IN query)  
GET /products?filter_category=1&filter_category=2&filter_category=3

# Custom filter logic in service
```

### Search Functionality

Configure searchable fields in your service:

```go
SearchFields: []string{"name", "description", "sku"}
```

Then use:
```bash
GET /products?search=gaming laptop
```

### Soft Deletes

Enable soft deletes in your configuration:

```go
config := &contracts.CrudConfig{
    SoftDelete: true,
    // ...
}
```

This provides automatic soft delete functionality with restore capabilities.

## Best Practices

1. **Interface Segregation**: Implement only the interfaces you need
2. **Progressive Enhancement**: Start with BaseCrudService, add features as needed
3. **Domain Logic**: Keep business logic in services, not controllers
4. **Validation**: Implement validation in services for reusability
5. **Error Handling**: Use consistent error responses via ResponseBuilder
6. **Configuration**: Use CrudConfig for customizable behavior

## Migration

Create migrations for your models:

```go
// database/migrations/create_products_table.go
func (receiver *CreateProductsTable) Up() error {
    return migration.Schema().Create("products", func(table schema.Blueprint) {
        table.ID()
        table.String("name").NotNull()
        table.Text("description").Nullable()
        table.Decimal("price", 10, 2).NotNull()
        table.Integer("stock").Default(0)
        table.UnsignedBigInteger("category_id").Nullable()
        table.Timestamps()
        table.SoftDeletes()
    })
}
```

## Testing

The modular design makes unit testing straightforward:

```go
func TestProductService_Create(t *testing.T) {
    service := services.NewProductService()
    
    data := map[string]interface{}{
        "name": "Test Product",
        "price": 99.99,
    }
    
    result, err := service.Create(data)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

This CRUD system provides a solid foundation for building consistent, scalable APIs while maintaining flexibility for domain-specific requirements.