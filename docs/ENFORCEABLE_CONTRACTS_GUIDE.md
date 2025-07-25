# Enforceable Contracts Guide

This guide explains how to use the new enforceable contract system in the Goravel blog application. The system uses Go interfaces and builder patterns to ensure compile-time checking of required methods and configurations.

## Problem Statement

Previously, our contracts (GenericCrudService, GenericCrudController, GenericPageController) were not enforceable. Users could:
- Skip implementing critical methods without compile-time errors
- Forget to call required initialization methods
- Miss essential configuration steps

## Solution: Interfaces + Builder Pattern

We've implemented a two-part solution:

1. **Go Interfaces**: Define contracts that must be implemented
2. **Builder Pattern**: Enforce required configuration steps at compile time

## Core Interfaces

### Service Interfaces

```go
// CrudServiceContract - Required for all CRUD services
type CrudServiceContract interface {
    // Core CRUD operations
    GetList(req ListRequest) (*PaginatedResult, error)
    GetByID(id uint) (interface{}, error)
    Create(data map[string]interface{}) (interface{}, error)
    Update(id uint, data map[string]interface{}) (interface{}, error)
    Delete(id uint) error
    
    // Search and filtering
    Search(req SearchRequest) (*PaginatedResult, error)
    GetListAdvanced(req ListRequest, filters map[string]interface{}) (*PaginatedResult, error)
    
    // Metadata
    GetSearchableFields() []string
    GetSortableFields() []string
    GetFilterableFields() []string
    GetValidationRules() map[string]interface{}
}
```

### Controller Interfaces

```go
// FullCrudControllerContract - Required for CRUD controllers
type FullCrudControllerContract interface {
    CrudControllerContract
    AuthorizationControllerContract
    SearchableControllerContract
}
```

### Page Controller Interface

```go
// PageControllerContract - Required for Inertia page controllers
type PageControllerContract interface {
    Index(ctx http.Context) http.Response
    CheckPermission(ctx http.Context, permission string, resource interface{}) error
    RequireAuthentication(ctx http.Context) error
    BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool
    ValidatePageRequest(ctx http.Context) (*ListRequest, error)
    GetProps(result *PaginatedResult, req *ListRequest, permissions map[string]bool, stats map[string]interface{}) PageProps
}
```

### Request Interfaces

```go
// CreateRequestContract - For create operations
type CreateRequestContract interface {
    RequestContract
    ToCreateData() map[string]interface{}
}

// UpdateRequestContract - For update operations
type UpdateRequestContract interface {
    RequestContract
    ToUpdateData() map[string]interface{}
    GetResourceID() interface{}
}
```

## Builder Patterns

### Service Builder

The service builder enforces that all required configurations are set:

```go
// Example: Creating a book service with enforced contracts
service := contracts.NewServiceBuilder[models.Book]("book", "id").
    WithSearchFields("title", "author", "isbn").        // REQUIRED
    WithSortFields("id", "title", "created_at").       // REQUIRED
    WithFilterFields("status", "author").               // REQUIRED
    WithValidationRules(map[string]interface{}{        // REQUIRED
        "title": "required|string|max:255",
        "author": "required|string|max:100",
    }).
    WithRelations("Tags", "Creator").                   // Optional
    WithSoftDeletes().                                  // Optional
    WithScopeFiltering("books", "created_by").          // Optional
    WithBeforeCreate(hookFunc).                         // Optional
    Build()

// If you skip any required step, you get a compile error!
// Try commenting out WithSearchFields() - the code won't compile!
```

### Controller Builder

The controller builder ensures all required methods are configured:

```go
// Example: Creating a book controller with enforced contracts
controller := contracts.NewControllerBuilder[models.Book, BookCreateRequest, BookUpdateRequest]("book").
    WithService(bookService).                           // REQUIRED
    WithAuthCheck(authCheckFunc).                       // REQUIRED
    WithRequestBindings(                                // REQUIRED
        bindCreateFunc,
        transformCreateFunc,
        bindUpdateFunc,
        transformUpdateFunc,
    ).
    WithBeforeStore(auditFunc).                         // Optional
    WithAfterStore(notifyFunc).                         // Optional
    Build()

// Skipping any required step results in compile error!
```

### Page Controller Builder

```go
// Example: Creating a page controller with enforced contracts
controller := contracts.NewPageControllerBuilder("books", "Books/BookList").
    WithService(bookService).                           // REQUIRED
    WithPagePermission(auth.PermissionRead).            // REQUIRED
    WithPropsConfiguration(propsFunc).                  // REQUIRED
    WithStatisticsSettings(true, statsFunc).            // Optional
    WithBeforeRender(hookFunc).                         // Optional
    Build()
```

### Request Builders

```go
// Create request with enforced validation
createRequest := contracts.NewCreateRequestBuilder().
    WithRules(map[string]string{                        // REQUIRED
        "title": "required|string|max:255",
        "author": "required|string|max:100",
    }).
    WithDataTransformer(transformFunc).                 // REQUIRED
    WithMessages(customMessages).                       // Optional
    WithAuthorization(authFunc).                        // Optional
    Build()

// Update request
updateRequest := contracts.NewUpdateRequestBuilder(resourceID).
    WithRules(validationRules).                         // REQUIRED
    WithDataTransformer(transformFunc).                 // REQUIRED
    Build()
```

## Migration Guide

### Step 1: Update Service

**Before (Non-enforceable):**
```go
type BookService struct {
    *contracts.GenericCrudService[models.Book]
}

func NewBookService() *BookService {
    service := &BookService{
        GenericCrudService: contracts.NewGenericCrudService[models.Book]("book"),
    }
    
    // Easy to forget these required calls!
    service.SetSearchFields("title", "author")
    service.SetSortFields("id", "title")
    // ... etc
    
    return service
}
```

**After (Enforceable):**
```go
func NewBookService() contracts.CrudServiceContract {
    return contracts.NewServiceBuilder[models.Book]("book", "id").
        WithSearchFields("title", "author", "isbn").    // Can't skip!
        WithSortFields("id", "title", "created_at").    // Can't skip!
        WithFilterFields("status", "author").            // Can't skip!
        WithValidationRules(rules).                      // Can't skip!
        Build()
}
```

### Step 2: Update Controller

**Before:**
```go
type BookController struct {
    *contracts.GenericCrudController[models.Book, BookCreateRequest, BookUpdateRequest]
}

func NewBookController() *BookController {
    // Could forget to set required fields!
    controller := &BookController{
        GenericCrudController: contracts.NewGenericCrudController[...](...),
    }
    // Easy to forget SetAuthCheck, SetRequestBindings, etc.
    return controller
}
```

**After:**
```go
func NewBookController() contracts.FullCrudControllerContract {
    return contracts.NewControllerBuilder[models.Book, BookCreateRequest, BookUpdateRequest]("book").
        WithService(bookService).                        // Can't skip!
        WithAuthCheck(authFunc).                         // Can't skip!
        WithRequestBindings(...).                        // Can't skip!
        Build()
}
```

## Benefits

1. **Compile-Time Safety**: Missing required configurations cause compile errors
2. **Self-Documenting**: The builder pattern shows exactly what's required vs optional
3. **Type Safety**: Go's type system ensures correct interface implementation
4. **No Runtime Surprises**: All contracts are verified at compile time

## Testing Contracts

Use the helper functions to ensure implementations:

```go
// These panic if the contract isn't implemented correctly
service := contracts.MustImplementCrudService(myService)
controller := contracts.MustImplementCrudController(myController)
pageController := contracts.MustImplementPageController(myPageController)
createRequest := contracts.MustImplementCreateRequest(myRequest)
```

## Best Practices

1. **Use Builders for New Code**: Always use the builder pattern for new services/controllers
2. **Migrate Gradually**: Update existing code one component at a time
3. **Test Interfaces**: Write tests against the interface, not the implementation
4. **Document Custom Methods**: If you extend beyond the base interface, document why

## Common Pitfalls to Avoid

1. **Don't Cast to Concrete Types**: Work with interfaces, not implementations
2. **Don't Skip Builders**: The old way still works but isn't enforceable
3. **Don't Ignore Compile Errors**: They're protecting you from runtime issues

## Conclusion

The new enforceable contract system provides compile-time guarantees that all required methods and configurations are properly implemented. This prevents runtime errors and makes the codebase more maintainable and reliable.