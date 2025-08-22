# Comprehensive CRUD End-to-End Guide

This guide demonstrates how to create a complete CRUD module/service in this framework using the Books module as an example. Follow these steps to create your own module from scratch.

## Table of Contents
1. [Model Creation](#1-model-creation)
2. [Database Migration](#2-database-migration)
3. [Permission Registration](#3-permission-registration)
4. [Service Layer](#4-service-layer)
5. [Request Validation](#5-request-validation)
6. [API Controller](#6-api-controller)
7. [Page Controller (Inertia)](#7-page-controller-inertia)
8. [Routing](#8-routing)
9. [Navigation Configuration](#9-navigation-configuration)
10. [Frontend Components](#10-frontend-components)

## 1. Model Creation

Create your model in `app/models/book.go`:

```go
package models

import (
    "encoding/json"
    "gorm.io/gorm"
    "time"
)

// Book entity - using BaseAuditableModel for consistent audit fields
type Book struct {
    BaseAuditableModel

    Title       string     `json:"title" gorm:"not null"`
    Author      string     `json:"author" gorm:"not null"`
    ISBN        string     `json:"isbn" gorm:"unique;not null"`
    Description string     `json:"description"`
    Price       float64    `json:"price" gorm:"default:0"`
    Status      string     `json:"status" gorm:"default:'AVAILABLE'"` 
    PublishedAt *time.Time `json:"publishedAt" gorm:"column:published_at"`
    TagsJSON    string     `json:"-" gorm:"column:tags;type:text"` // Store as JSON in DB
    Tags        []string   `json:"tags" gorm:"-"`                  // Virtual field for API
}

// BeforeSave hook to convert tags array to JSON
func (b *Book) BeforeSave(tx *gorm.DB) error {
    if len(b.Tags) > 0 {
        tagsBytes, err := json.Marshal(b.Tags)
        if err != nil {
            return err
        }
        b.TagsJSON = string(tagsBytes)
    } else {
        b.TagsJSON = ""
    }
    return nil
}

// AfterFind hook to convert tags JSON to array
func (b *Book) AfterFind(tx *gorm.DB) error {
    if b.TagsJSON != "" {
        err := json.Unmarshal([]byte(b.TagsJSON), &b.Tags)
        if err != nil {
            b.Tags = []string{}
        }
    } else {
        b.Tags = []string{}
    }
    return nil
}

// SearchFields returns the fields that can be searched
func (b Book) SearchFields() []string {
    return []string{"title", "author", "isbn", "description"}
}

// TableName returns the table name for this model
func (b Book) TableName() string {
    return "books"
}

// MarshalJSON custom JSON marshaling to handle date formatting
func (b Book) MarshalJSON() ([]byte, error) {
    type Alias Book
    return json.Marshal(&struct {
        *Alias
        CreatedAt   *time.Time `json:"createdAt,omitempty"`
        UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
        PublishedAt *time.Time `json:"publishedAt,omitempty"`
    }{
        Alias:       (*Alias)(&b),
        CreatedAt:   timeToPtr(b.CreatedAt.StdTime()),
        UpdatedAt:   timeToPtr(b.UpdatedAt.StdTime()),
        PublishedAt: b.PublishedAt,
    })
}
```

### Key Model Features:
- **BaseAuditableModel**: Provides ID, CreatedAt, UpdatedAt, DeletedAt, CreatedBy, UpdatedBy fields
- **GORM Hooks**: BeforeSave/AfterFind for data transformation
- **JSON Handling**: Custom marshaling for consistent date formats
- **Searchable Interface**: Define which fields can be searched

## 2. Database Migration

Create migration in `database/migrations/YYYYMMDDHHMMSS_create_books_table.go`:

```go
package migrations

import (
    "gorm.io/gorm"
)

type CreateBooksTable struct{}

func (m *CreateBooksTable) Up() func(*gorm.DB) error {
    return func(tx *gorm.DB) error {
        type Book struct {
            gorm.Model
            Title       string  `gorm:"not null;index"`
            Author      string  `gorm:"not null;index"`
            ISBN        string  `gorm:"unique;not null"`
            Description string  `gorm:"type:text"`
            Price       float64 `gorm:"default:0"`
            Status      string  `gorm:"default:'AVAILABLE';index"`
            PublishedAt *time.Time
            Tags        string  `gorm:"type:text"`
            CreatedBy   uint    `gorm:"index"`
            UpdatedBy   uint
        }
        
        return tx.AutoMigrate(&Book{})
    }
}

func (m *CreateBooksTable) Down() func(*gorm.DB) error {
    return func(tx *gorm.DB) error {
        return tx.Migrator().DropTable("books")
    }
}

// Signature returns the migration signature
func (m *CreateBooksTable) Signature() string {
    return "2024_01_01_000000_create_books_table"
}
```

## 3. Permission Registration

Add your service to the permission registry in `app/auth/permission_helper.go`:

```go
// In app/auth/permissions.go or similar file

// Add to ServiceRegistry
const (
    ServiceBooks ServiceRegistry = "books"
)

// Register in GetAllServices()
func GetAllServices() []ServiceDefinition {
    return []ServiceDefinition{
        // ... existing services
        {
            Service:     ServiceBooks,
            Name:        "Books Management",
            Description: "Manage library books",
        },
    }
}
```

The framework automatically generates permissions for your service:
- `books_view` / `books_view_by_all` / `books_view_by_my_role` / `books_view_by_me`
- `books_read` / `books_read_by_all` / `books_read_by_my_role` / `books_read_by_me`
- `books_create`
- `books_update` / `books_update_by_all` / `books_update_by_my_role` / `books_update_by_me`
- `books_delete` / `books_delete_by_all` / `books_delete_by_my_role` / `books_delete_by_me`
- `books_export`
- `books_bulk_update`
- `books_bulk_delete`
- `books_manage`

## 4. Service Layer

Create your service in `app/services/book_service.go`:

```go
package services

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/goravel/framework/facades"
    "players/app/contracts"
    "players/app/models"
)

// BookService implements book-specific business logic
type BookService struct {
    baseService contracts.CrudServiceContract
}

// Implement CrudServiceContract interface by delegating to baseService
func (s *BookService) GetList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
    return s.baseService.GetList(req)
}

func (s *BookService) GetByID(id uint) (interface{}, error) {
    return s.baseService.GetByID(id)
}

func (s *BookService) Create(data map[string]interface{}) (interface{}, error) {
    return s.baseService.Create(data)
}

func (s *BookService) Update(id uint, data map[string]interface{}) (interface{}, error) {
    return s.baseService.Update(id, data)
}

func (s *BookService) Delete(id uint) error {
    return s.baseService.Delete(id)
}

// NewBookService creates a new book service using the builder pattern
func NewBookService() *BookService {
    service := contracts.NewServiceBuilder[models.Book]("books", "id").
        WithSearchFields("title", "author", "isbn", "description").    // REQUIRED
        WithSortFields("id", "title", "author", "price", "status",    // REQUIRED
            "created_at", "updated_at", "published_at").
        WithFilterFields("status", "author").                          // REQUIRED
        WithValidationRules(map[string]interface{}{                    // REQUIRED
            "title":       "required|string|max:255",
            "author":      "required|string|max:100",
            "isbn":        "required|string|max:20",
            "status":      "required|string|in:AVAILABLE,BORROWED,MAINTENANCE,RESERVED",
            "price":       "numeric|min:0",
            "publishedAt": "date",
            "description": "string|max:1000",
            "tags":        "array",
            "tags.*":      "string|max:50",
        }).
        WithRelations("Creator", "Updater").                          // Optional
        WithDefaultSort("created_at", "DESC").                        // Optional
        WithSoftDeletes().                                           // Optional
        WithScopeFiltering("books", "created_by").                   // Optional
        WithBeforeCreate(func(data map[string]interface{}) error {   // Optional
            // Handle tags array to JSON conversion
            if tags, exists := data["tags"]; exists && tags != nil {
                // Convert tags to JSON string for storage
                // ... conversion logic
            }
            return nil
        }).
        WithBeforeUpdate(func(id uint, data map[string]interface{}) error { // Optional
            // Similar tags handling for updates
            return nil
        }).
        Build()

    bookServiceInstance := &BookService{
        baseService: service,
    }

    // Set the actual service instance for proper method resolution
    contracts.SetActualServiceHelper(service, bookServiceInstance, "BookService")

    return bookServiceInstance
}

// Override methods for custom behavior
func (s *BookService) GetColumnMapping() map[string]string {
    mapping := s.baseService.GetColumnMapping()
    // Add book-specific mappings
    mapping["publishedAt"] = "published_at"
    mapping["createdAt"] = "created_at"
    mapping["updatedAt"] = "updated_at"
    return mapping
}

// GetFilterDefinitions returns filter definitions for the books resource
func (s *BookService) GetFilterDefinitions() []contracts.FilterDefinition {
    return []contracts.FilterDefinition{
        // String filters - automatically get all string operators
        contracts.NewFilterDefinition(
            "title",
            "Title",
            contracts.FilterTypeString,
            nil, // Will use GetOperatorsForType(FilterTypeString)
        ),
        // Number filters - automatically get all number operators
        contracts.NewFilterDefinition(
            "price",
            "Price",
            contracts.FilterTypeNumber,
            nil, // Will use GetOperatorsForType(FilterTypeNumber)
        ),
        // Enum filter
        {
            Field: "status",
            Label: "Status",
            Type:  contracts.FilterTypeEnum,
            Operators: nil, // Will use GetOperatorsForType(FilterTypeEnum)
            EnumValues: []string{
                "AVAILABLE",
                "BORROWED",
                "MAINTENANCE",
                "RESERVED",
            },
        },
        // Date filter
        contracts.NewFilterDefinition(
            "published_at",
            "Published Date",
            contracts.FilterTypeDate,
            nil, // Will use GetOperatorsForType(FilterTypeDate)
        ),
    }
}

// Add custom business logic methods
func (s *BookService) GetByISBN(isbn string) (*models.Book, error) {
    var book models.Book
    err := facades.Orm().Query().
        Model(&models.Book{}).
        Where("isbn = ?", isbn).
        With("Creator").
        With("Updater").
        First(&book)

    if err != nil {
        return nil, err
    }

    if book.ID == 0 {
        return nil, fmt.Errorf("book not found")
    }

    return &book, nil
}

// BorrowBook updates book status to borrowed
func (s *BookService) BorrowBook(id uint) error {
    bookInterface, err := s.GetByID(id)
    if err != nil {
        return err
    }

    book, ok := bookInterface.(*models.Book)
    if !ok {
        return errors.New("invalid book type")
    }

    if book.Status != "AVAILABLE" {
        return errors.New("book is not available for borrowing")
    }

    updateData := map[string]interface{}{
        "status": "BORROWED",
    }

    _, err = s.Update(id, updateData)
    return err
}

// GetBookStatistics returns statistics about books
func (s *BookService) GetBookStatistics() (map[string]interface{}, error) {
    var stats struct {
        TotalBooks       int64
        AvailableBooks   int64
        BorrowedBooks    int64
        MaintenanceBooks int64
        TotalValue       float64
        AveragePrice     float64
    }

    // Get total books (excluding soft deleted)
    facades.Orm().Query().Model(&models.Book{}).
        Where("deleted_at IS NULL").Count(&stats.TotalBooks)

    // Get available books
    facades.Orm().Query().Model(&models.Book{}).
        Where("status = ? AND deleted_at IS NULL", "AVAILABLE").
        Count(&stats.AvailableBooks)

    // ... more statistics queries

    return map[string]interface{}{
        "totalBooks":       stats.TotalBooks,
        "availableBooks":   stats.AvailableBooks,
        "borrowedBooks":    stats.BorrowedBooks,
        "maintenanceBooks": stats.MaintenanceBooks,
        "totalValue":       fmt.Sprintf("$%.2f", stats.TotalValue),
        "averagePrice":     stats.AveragePrice,
    }, nil
}
```

### Service Builder Features:
- **Required Methods**: WithSearchFields, WithSortFields, WithFilterFields, WithValidationRules
- **Optional Methods**: WithRelations, WithDefaultSort, WithSoftDeletes, WithScopeFiltering, WithBeforeCreate/Update/Delete
- **Base CRUD Operations**: GetList, GetByID, Create, Update, Delete, Search, GetListAdvanced
- **Filter Support**: Automatic operator mapping based on field type
- **Column Mapping**: Frontend to database field mapping
- **Custom Business Logic**: Add your own methods like BorrowBook, GetByISBN, etc.

## 5. Request Validation

Create request validators in `app/http/requests/book_request.go`:

```go
package requests

import (
    "fmt"
    "players/app/contracts"
    "strings"
    "github.com/goravel/framework/contracts/http"
)

// BookCreateRequest handles book creation validation
type BookCreateRequest struct {
    Title       string   `form:"title" json:"title"`
    Author      string   `form:"author" json:"author"`
    ISBN        string   `form:"isbn" json:"isbn"`
    Description string   `form:"description" json:"description"`
    Price       float64  `form:"price" json:"price"`
    Status      string   `form:"status" json:"status"`
    PublishedAt string   `form:"publishedAt" json:"publishedAt"`
    Tags        []string `form:"tags" json:"tags"`
}

// Rules defines validation rules
func (r *BookCreateRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "title":  "required",
        "author": "required",
        "isbn":   "required",
        "price":  fmt.Sprintf("required|%s", fmt.Sprintf(contracts.MinValue, 0)),
        "status": "in:AVAILABLE,BORROWED,MAINTENANCE,RESERVED",
    }
}

// Messages defines custom validation messages
func (r *BookCreateRequest) Messages(ctx http.Context) map[string]string {
    return map[string]string{
        "title.required":  "Book title is required",
        "author.required": "Author name is required",
        "isbn.required":   "ISBN is required",
        "price.required":  "Price is required",
        "price.min":       "Price must be greater than or equal to 0",
        "status.in":       "Status must be one of: AVAILABLE, BORROWED, MAINTENANCE, RESERVED",
    }
}

// PrepareForValidation allows modification of input before validation
func (r *BookCreateRequest) PrepareForValidation(ctx http.Context) error {
    // Normalize ISBN by removing hyphens
    if r.ISBN != "" {
        r.ISBN = strings.ReplaceAll(r.ISBN, "-", "")
        r.ISBN = strings.ReplaceAll(r.ISBN, " ", "")
    }

    // Set default status if not provided
    if r.Status == "" {
        r.Status = "AVAILABLE"
    }

    return nil
}

// ToCreateData converts the request to create data map
func (r *BookCreateRequest) ToCreateData() map[string]interface{} {
    data := map[string]interface{}{
        "title":       r.Title,
        "author":      r.Author,
        "isbn":        r.ISBN,
        "description": r.Description,
        "price":       r.Price,
        "status":      r.Status,
    }

    if r.PublishedAt != "" {
        data["publishedAt"] = r.PublishedAt
    }

    if len(r.Tags) > 0 {
        data["tags"] = r.Tags
    }

    return data
}

// BookUpdateRequest handles book update validation
type BookUpdateRequest struct {
    Title       *string   `form:"title" json:"title"`
    Author      *string   `form:"author" json:"author"`
    ISBN        *string   `form:"isbn" json:"isbn"`
    Description *string   `form:"description" json:"description"`
    Price       *float64  `form:"price" json:"price"`
    Status      *string   `form:"status" json:"status"`
    PublishedAt *string   `form:"publishedAt" json:"publishedAt"`
    Tags        *[]string `form:"tags" json:"tags"`
    ID          uint      `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for updates
func (r *BookUpdateRequest) Rules(ctx http.Context) map[string]string {
    rules := map[string]string{}

    // Only validate fields that are provided
    if r.Title != nil {
        rules["title"] = "required"
    }
    if r.Author != nil {
        rules["author"] = "required"
    }
    if r.Price != nil {
        rules["price"] = fmt.Sprintf("required|%s", fmt.Sprintf(contracts.MinValue, 0))
    }
    if r.Status != nil {
        rules["status"] = "in:AVAILABLE,BORROWED,MAINTENANCE,RESERVED"
    }

    // Prevent empty rules error
    if len(rules) == 0 {
        rules["_at_least_one_field"] = "sometimes"
    }

    return rules
}

// ToUpdateData converts the request to update data map
func (r *BookUpdateRequest) ToUpdateData() map[string]interface{} {
    data := map[string]interface{}{}

    // Only include fields that are provided (not nil)
    if r.Title != nil {
        data["title"] = *r.Title
    }
    if r.Author != nil {
        data["author"] = *r.Author
    }
    if r.ISBN != nil {
        data["isbn"] = *r.ISBN
    }
    // ... other fields

    return data
}
```

## 6. API Controller

Create your API controller in `app/http/controllers/books/book_controller.go`:

```go
package books

import (
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
    "players/app/auth"
    "players/app/contracts"
    "players/app/http/requests"
    "players/app/models"
    "players/app/services"
)

// BookController handles API endpoints for book management
type BookController struct {
    *contracts.StaticEnforcedController[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest]
    bookService *services.BookService
}

// NewBookController creates a new book controller with compile-time enforcement
func NewBookController() *BookController {
    bookService := services.NewBookService()

    // Build controller with compile-time enforcement
    staticController := contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
        "book",
        bookService,
    ).
        ValidateCreateRequest().
        ValidateUpdateRequest().
        WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
            scopedHelper := auth.GetScopedPermissionHelper()

            // Map generic actions to permission actions
            var permAction auth.CorePermissionAction
            switch action {
            case "viewAny", "view":
                permAction = auth.PermissionRead
            case "create":
                permAction = auth.PermissionCreate
            case "update":
                permAction = auth.PermissionUpdate
            case "delete":
                permAction = auth.PermissionDelete
            default:
                permAction = auth.PermissionManage
            }

            // For specific resource actions, pass the resource
            if resource != nil {
                _, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBooks, permAction, resource)
                return err
            }

            // For general actions, check without specific resource
            _, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBooks, permAction, nil)
            return err
        }).
        Build()

    controller := &BookController{
        StaticEnforcedController: staticController,
        bookService:              bookService,
    }

    // Set custom hooks
    controller.SetBeforeStore(func(ctx http.Context, data map[string]interface{}) error {
        // Set created_by from authenticated user
        var user models.User
        if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
            data["created_by"] = user.ID
        }
        return nil
    })

    return controller
}

// The controller automatically implements these methods from StaticEnforcedController:
// - Index (GET /api/books)
// - Show (GET /api/books/{id})
// - Store (POST /api/books)
// - Update (PUT /api/books/{id})
// - Delete (DELETE /api/books/{id})
// - Search (GET /api/books/search)
// - FilterMetadata (GET /api/books/filters)

// Add custom endpoints
func (c *BookController) Borrow(ctx http.Context) http.Response {
    id, err := c.ValidateID(ctx, "id")
    if err != nil {
        return c.BadRequestResponse(ctx, "Invalid book ID", nil)
    }

    // Check permissions
    if err := c.CheckAuth(ctx, "update", nil); err != nil {
        return c.ForbiddenResponse(ctx, "Access denied")
    }

    // Borrow the book
    err = c.bookService.BorrowBook(id)
    if err != nil {
        return c.BadRequestResponse(ctx, err.Error(), nil)
    }

    return c.SuccessResponse(ctx, map[string]interface{}{
        "message": "Book borrowed successfully",
    }, "Book borrowed")
}

// Statistics endpoint
func (c *BookController) Statistics(ctx http.Context) http.Response {
    // Check permissions
    if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
        return c.ForbiddenResponse(ctx, "Access denied")
    }

    stats, err := c.bookService.GetBookStatistics()
    if err != nil {
        return c.InternalErrorResponse(ctx, "Failed to retrieve statistics")
    }

    return c.SuccessResponse(ctx, stats, "Statistics retrieved successfully")
}

// Override GetFilterDefinitions to provide custom filters
func (c *BookController) GetFilterDefinitions() []contracts.FilterDefinition {
    return []contracts.FilterDefinition{
        contracts.NewFilterDefinition(
            "price",
            "Price",
            contracts.FilterTypeNumber,
            nil, // Automatically uses all number operators
        ),
        {
            Field: "status",
            Label: "Status",
            Type:  contracts.FilterTypeEnum,
            Operators: nil, // Automatically uses enum operators
            EnumValues: []string{"AVAILABLE", "BORROWED", "MAINTENANCE", "RESERVED"},
        },
        // ... more filters
    }
}
```

### Controller Features:
- **Compile-time Safety**: Uses generics to ensure type safety
- **Automatic CRUD**: Index, Show, Store, Update, Delete methods
- **Permission Checking**: Integrated auth checking with scoped permissions
- **Request Validation**: Automatic validation using request objects
- **Custom Hooks**: BeforeStore, BeforeUpdate, AfterStore, etc.
- **Response Helpers**: SuccessResponse, BadRequestResponse, ForbiddenResponse, etc.

## 7. Page Controller (Inertia)

Create the Inertia page controller in `app/http/controllers/books/books_page_controller.go`:

```go
package books

import (
    "players/app/auth"
    "players/app/contracts"
    "players/app/services"
)

// BooksPageController handles the books page
type BooksPageController struct {
    *contracts.GenericPageController
    bookService *services.BookService
}

// NewBooksPageController creates a new books page controller
func NewBooksPageController() *BooksPageController {
    bookService := services.NewBookService()

    return &BooksPageController{
        GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
            ResourceType:      "books",
            PageComponent:     "Books/Index",          // React component path
            Service:           bookService,
            ServiceIdentifier: auth.ServiceBooks,
            StatsEnabled:      true,
            StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
                stats, _ := bookService.GetBookStatistics()
                return stats
            },
        }),
        bookService: bookService,
    }
}
```

### Page Controller Features:
- **Automatic Props**: Permissions, filters, sorting, pagination config
- **Stats Support**: Optional statistics data for dashboards
- **Service Integration**: Connects to your service for data
- **Inertia Rendering**: Handles server-side rendering setup

## 8. Routing

### API Routes (`routes/api.go`)

```go
// In Api() function

bookController := books.NewBookController()

// Book resource routes (with optional auth for scoped permissions)
router.Middleware(optionalAuth).Group(func(optionalAuthRouter route.Router) {
    optionalAuthRouter.Get("/books", bookController.Index)
    optionalAuthRouter.Get("/books/search", bookController.Search)
    optionalAuthRouter.Get("/books/filters", bookController.FilterMetadata)
    optionalAuthRouter.Get("/books/available", bookController.Available)
    optionalAuthRouter.Get("/books/isbn/{isbn}", bookController.GetByISBN)
    optionalAuthRouter.Get("/books/author/{author}", bookController.GetByAuthor)
    optionalAuthRouter.Get("/books/{id}", bookController.Show) // Must be last
})

// Protected routes (require authentication)
router.Middleware(jwtAuth).Group(func(protectedRouter route.Router) {
    // Book CRUD routes
    protectedRouter.Post("/books", bookController.Store)
    protectedRouter.Put("/books/{id}", bookController.Update)
    protectedRouter.Delete("/books/{id}", bookController.Delete)
    
    // Custom endpoints
    protectedRouter.Post("/books/{id}/borrow", bookController.Borrow)
    protectedRouter.Post("/books/{id}/return", bookController.Return)
    protectedRouter.Get("/books/statistics", bookController.Statistics)
})
```

### Web Routes (`routes/web.go`)

```go
// In Web() function

booksPageController := books.NewBooksPageController()

// Inside authenticated routes group
router.Get("/admin/books", booksPageController.Index)
```

## 9. Navigation Configuration

Add your module to `resources/js/config/navigation.ts`:

```typescript
import { BookIcon } from "lucide-react"

export const navigationConfig: NavigationConfig = {
    navMain: [
        {
            title: "Books",
            url: "/admin/books",
            icon: BookIcon,
            requiredService: "books",
            requiredAction: "read" as const,
        },
        // ... other nav items
    ],
    // ... other navigation sections
}
```

## 10. Frontend Components

### Component Structure
```
resources/js/pages/Books/
├── Index.tsx              # Main page component
└── sections/
    ├── index.ts          # Export barrel
    ├── bookPageConfig.tsx # Page configuration
    ├── BookColumns.tsx    # Table column definitions
    ├── BookCreateForm.tsx # Create form component
    ├── BookEditForm.tsx   # Edit form component
    └── BookDetailView.tsx # Detail view component
```

### Main Page Component (`Index.tsx`)

```tsx
import React from "react";
import { CrudPage } from "@/components/Crud/CrudPage";
import { bookPageConfig } from "./sections/bookPageConfig";

export default function BooksIndex() {
    return <CrudPage {...bookPageConfig} />;
}
```

### Page Configuration (`bookPageConfig.tsx`)

```tsx
import { CrudPageConfig } from "@/components/Crud/types";
import { BookColumns } from "./BookColumns";
import { BookCreateForm } from "./BookCreateForm";
import { BookEditForm } from "./BookEditForm";
import { BookDetailView } from "./BookDetailView";

export const bookPageConfig: CrudPageConfig = {
    // API configuration
    apiEndpoint: "/api/books",
    resourceName: "book",
    resourceNamePlural: "books",
    
    // Component configuration
    columns: BookColumns,
    createForm: BookCreateForm,
    editForm: BookEditForm,
    detailView: BookDetailView,
    
    // Feature flags
    features: {
        create: true,
        edit: true,
        delete: true,
        search: true,
        filters: true,
        export: true,
        bulkActions: true,
    },
    
    // UI configuration
    ui: {
        createButtonText: "Add New Book",
        pageTitle: "Books Management",
        pageDescription: "Manage your library book collection",
        emptyStateMessage: "No books found. Add your first book to get started.",
    },
    
    // Default settings
    defaultSort: {
        field: "created_at",
        direction: "desc",
    },
    
    // Stats configuration (optional)
    stats: {
        enabled: true,
        endpoint: "/api/books/statistics",
    },
};
```

### Table Columns (`BookColumns.tsx`)

```tsx
import { ColumnDef } from "@tanstack/react-table";
import { Badge } from "@/components/ui/badge";
import { formatDate, formatCurrency } from "@/lib/utils";
import { DataTableColumnHeader } from "@/components/data-table/DataTableColumnHeader";

export const BookColumns: ColumnDef<any>[] = [
    {
        accessorKey: "title",
        header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Title" />
        ),
        cell: ({ row }) => {
            return (
                <div className="flex flex-col">
                    <span className="font-medium">{row.getValue("title")}</span>
                    <span className="text-sm text-muted-foreground">
                        ISBN: {row.original.isbn}
                    </span>
                </div>
            );
        },
    },
    {
        accessorKey: "author",
        header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Author" />
        ),
    },
    {
        accessorKey: "status",
        header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Status" />
        ),
        cell: ({ row }) => {
            const status = row.getValue("status") as string;
            const statusConfig = {
                AVAILABLE: { label: "Available", variant: "success" },
                BORROWED: { label: "Borrowed", variant: "warning" },
                MAINTENANCE: { label: "Maintenance", variant: "secondary" },
                RESERVED: { label: "Reserved", variant: "info" },
            };
            
            const config = statusConfig[status] || { 
                label: status, 
                variant: "default" 
            };
            
            return (
                <Badge variant={config.variant as any}>
                    {config.label}
                </Badge>
            );
        },
        filterFn: (row, id, value) => {
            return value.includes(row.getValue(id));
        },
    },
    {
        accessorKey: "price",
        header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Price" />
        ),
        cell: ({ row }) => formatCurrency(row.getValue("price")),
    },
    {
        accessorKey: "publishedAt",
        header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Published" />
        ),
        cell: ({ row }) => {
            const date = row.getValue("publishedAt");
            return date ? formatDate(date as string) : "-";
        },
    },
    {
        accessorKey: "tags",
        header: "Tags",
        cell: ({ row }) => {
            const tags = row.getValue("tags") as string[];
            if (!tags || tags.length === 0) return "-";
            
            return (
                <div className="flex flex-wrap gap-1">
                    {tags.map((tag, index) => (
                        <Badge key={index} variant="outline" className="text-xs">
                            {tag}
                        </Badge>
                    ))}
                </div>
            );
        },
    },
];
```

### Create Form (`BookCreateForm.tsx`)

```tsx
import React from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { TagInput } from "@/components/ui/tag-input";
import { DatePicker } from "@/components/ui/date-picker";

const bookSchema = z.object({
    title: z.string().min(1, "Title is required").max(255),
    author: z.string().min(1, "Author is required").max(100),
    isbn: z.string().min(1, "ISBN is required").max(20),
    description: z.string().max(1000).optional(),
    price: z.number().min(0, "Price must be positive"),
    status: z.enum(["AVAILABLE", "BORROWED", "MAINTENANCE", "RESERVED"]),
    publishedAt: z.string().optional(),
    tags: z.array(z.string()).optional(),
});

type BookFormData = z.infer<typeof bookSchema>;

interface BookCreateFormProps {
    onSubmit: (data: any) => Promise<void>;
    onCancel: () => void;
    isSubmitting?: boolean;
}

export function BookCreateForm({ 
    onSubmit, 
    onCancel, 
    isSubmitting 
}: BookCreateFormProps) {
    const form = useForm<BookFormData>({
        resolver: zodResolver(bookSchema),
        defaultValues: {
            title: "",
            author: "",
            isbn: "",
            description: "",
            price: 0,
            status: "AVAILABLE",
            tags: [],
        },
    });

    const handleSubmit = async (data: BookFormData) => {
        await onSubmit(data);
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
                <div className="grid grid-cols-2 gap-4">
                    <FormField
                        control={form.control}
                        name="title"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Title</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter book title" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="author"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Author</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter author name" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="isbn"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>ISBN</FormLabel>
                                <FormControl>
                                    <Input placeholder="Enter ISBN" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="price"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Price</FormLabel>
                                <FormControl>
                                    <Input 
                                        type="number" 
                                        step="0.01"
                                        placeholder="0.00" 
                                        {...field}
                                        onChange={(e) => field.onChange(parseFloat(e.target.value))}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="status"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Status</FormLabel>
                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                    <FormControl>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select status" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent>
                                        <SelectItem value="AVAILABLE">Available</SelectItem>
                                        <SelectItem value="BORROWED">Borrowed</SelectItem>
                                        <SelectItem value="MAINTENANCE">Maintenance</SelectItem>
                                        <SelectItem value="RESERVED">Reserved</SelectItem>
                                    </SelectContent>
                                </Select>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="publishedAt"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Published Date</FormLabel>
                                <FormControl>
                                    <DatePicker
                                        value={field.value}
                                        onChange={field.onChange}
                                        placeholder="Select date"
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                </div>

                <FormField
                    control={form.control}
                    name="description"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Description</FormLabel>
                            <FormControl>
                                <Textarea 
                                    placeholder="Enter book description"
                                    rows={3}
                                    {...field}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <FormField
                    control={form.control}
                    name="tags"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Tags</FormLabel>
                            <FormControl>
                                <TagInput
                                    value={field.value || []}
                                    onChange={field.onChange}
                                    placeholder="Add tags..."
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <div className="flex justify-end space-x-2">
                    <Button
                        type="button"
                        variant="outline"
                        onClick={onCancel}
                        disabled={isSubmitting}
                    >
                        Cancel
                    </Button>
                    <Button type="submit" disabled={isSubmitting}>
                        {isSubmitting ? "Creating..." : "Create Book"}
                    </Button>
                </div>
            </form>
        </Form>
    );
}
```

### Edit Form (`BookEditForm.tsx`)

```tsx
// Similar to BookCreateForm but with:
// 1. Initial data population from props
// 2. Optional fields (using Partial<BookFormData>)
// 3. Different submit button text
// 4. Pre-populated form values

export function BookEditForm({ 
    data,
    onSubmit, 
    onCancel, 
    isSubmitting 
}: BookEditFormProps) {
    const form = useForm<BookFormData>({
        resolver: zodResolver(bookSchema),
        defaultValues: {
            title: data.title || "",
            author: data.author || "",
            isbn: data.isbn || "",
            description: data.description || "",
            price: data.price || 0,
            status: data.status || "AVAILABLE",
            publishedAt: data.publishedAt,
            tags: data.tags || [],
        },
    });
    
    // ... rest similar to create form
}
```

### Detail View (`BookDetailView.tsx`)

```tsx
import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { formatDate, formatCurrency } from "@/lib/utils";

interface BookDetailViewProps {
    data: any;
}

export function BookDetailView({ data }: BookDetailViewProps) {
    return (
        <div className="space-y-6">
            <Card>
                <CardHeader>
                    <CardTitle>Book Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Title
                            </label>
                            <p className="mt-1">{data.title}</p>
                        </div>
                        
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Author
                            </label>
                            <p className="mt-1">{data.author}</p>
                        </div>
                        
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                ISBN
                            </label>
                            <p className="mt-1">{data.isbn}</p>
                        </div>
                        
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Price
                            </label>
                            <p className="mt-1">{formatCurrency(data.price)}</p>
                        </div>
                        
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Status
                            </label>
                            <div className="mt-1">
                                <Badge>{data.status}</Badge>
                            </div>
                        </div>
                        
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Published Date
                            </label>
                            <p className="mt-1">
                                {data.publishedAt ? formatDate(data.publishedAt) : "-"}
                            </p>
                        </div>
                    </div>
                    
                    {data.description && (
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Description
                            </label>
                            <p className="mt-1">{data.description}</p>
                        </div>
                    )}
                    
                    {data.tags && data.tags.length > 0 && (
                        <div>
                            <label className="text-sm font-medium text-muted-foreground">
                                Tags
                            </label>
                            <div className="mt-1 flex flex-wrap gap-1">
                                {data.tags.map((tag: string, index: number) => (
                                    <Badge key={index} variant="outline">
                                        {tag}
                                    </Badge>
                                ))}
                            </div>
                        </div>
                    )}
                </CardContent>
            </Card>
            
            {/* Audit Information */}
            <Card>
                <CardHeader>
                    <CardTitle>Audit Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                    <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">Created By</span>
                        <span className="text-sm">
                            {data.creator?.name || "System"}
                        </span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">Created At</span>
                        <span className="text-sm">{formatDate(data.createdAt)}</span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">Last Updated By</span>
                        <span className="text-sm">
                            {data.updater?.name || "System"}
                        </span>
                    </div>
                    <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">Last Updated At</span>
                        <span className="text-sm">{formatDate(data.updatedAt)}</span>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
```

## Key Features of the Framework

### 1. **Service Builder Pattern**
- Fluent API for configuring services
- Automatic CRUD operations
- Built-in search, filter, sort, and pagination
- Hooks for custom business logic

### 2. **Scoped Permissions**
- Row-level security with `_by_me`, `_by_my_role`, `_by_all` suffixes
- Automatic permission generation per service
- Permission checking at controller and service levels

### 3. **Filter System**
- Type-based operator mapping
- Support for string, number, date, boolean, enum, and array filters
- Compound filters with AND/OR logic
- Custom filter definitions per resource

### 4. **Frontend Integration**
- Reusable CRUD components
- Automatic form generation
- Built-in validation
- Responsive data tables with sorting, filtering, and pagination
- Real-time updates via WebSocket/SSE

### 5. **API Standards**
- RESTful endpoints
- Consistent response formats
- Automatic validation error handling
- Pagination metadata
- Filter metadata endpoints

## Best Practices

1. **Always use the service builder** for consistency
2. **Define proper validation rules** in both service and request objects
3. **Implement scoped permissions** for row-level security
4. **Use TypeScript interfaces** for frontend type safety
5. **Follow the naming conventions** for files and components
6. **Add proper indexes** to your database migrations
7. **Write custom business logic** in service methods, not controllers
8. **Use hooks** for data transformation and side effects
9. **Implement proper error handling** with meaningful messages
10. **Document your custom endpoints** and business rules

## Testing

### Unit Tests
- Test service methods independently
- Mock database operations
- Test validation rules

### Integration Tests
- Test full API endpoints
- Include authentication/authorization
- Test filter and search functionality

### Frontend Tests
- Component testing with React Testing Library
- Form validation tests
- Table interaction tests

## Common Patterns

### Adding a Status Workflow
```go
// In service
func (s *BookService) ChangeStatus(id uint, newStatus string) error {
    // Validate status transition
    // Update status
    // Log status change
    // Send notifications
}
```

### Adding Bulk Operations
```go
// In controller
func (c *BookController) BulkUpdateStatus(ctx http.Context) http.Response {
    var request struct {
        IDs    []uint `json:"ids"`
        Status string `json:"status"`
    }
    // Validate and process
}
```

### Adding Export Functionality
```go
// In service
func (s *BookService) ExportToCSV(filters map[string]interface{}) ([]byte, error) {
    // Query with filters
    // Generate CSV
    // Return bytes
}
```

This guide provides a complete blueprint for creating CRUD modules in the framework. Each module follows the same pattern, making it easy to maintain consistency across your application.