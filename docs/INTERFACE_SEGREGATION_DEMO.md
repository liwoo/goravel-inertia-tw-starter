# Interface Segregation for Compile-Time Security Enforcement

This document demonstrates how interface segregation enforces security requirements at compile time.

## The Problem

Without compile-time enforcement, developers can forget to add auth checking or validation:

```go
// This compiles but is insecure!
controller := SomeController{service: bookService}
route.Get("/api/books", controller.Index) // No auth check!
```

## The Solution: Interface Segregation

We use Go's type system to make insecure code impossible to compile.

### 1. Required Security Interfaces

```go
// Controllers MUST implement these interfaces to be registrable
type SecuredCrudController interface {
    AuthRequiredController      // Must have GetAuthChecker()
    ValidatedCreateController   // Must have GetCreateValidator()
    ValidatedUpdateController   // Must have GetUpdateValidator()
}
```

### 2. Secure Route Registration

```go
// This function ONLY accepts controllers that implement SecuredCrudController
func (r *SecureRouteRegistrar) RegisterSecuredCRUD(
    path string,
    controller SecuredCrudController, // <-- Compile-time enforcement!
    crudController CrudControllerContract,
) {
    // Registration code...
}
```

### 3. What Happens When You Try to Skip Security?

#### Example 1: Controller Without Auth Checker

```go
type InsecureController struct {
    service CrudServiceContract
}

// Missing: GetAuthChecker() method
// Missing: GetCreateValidator() method  
// Missing: GetUpdateValidator() method

// Try to register it:
registrar := NewSecureRouteRegistrar(route)
registrar.RegisterSecuredCRUD("/api/items", controller, controller)

// COMPILE ERROR:
// cannot use controller (variable of type *InsecureController) as SecuredCrudController value:
// *InsecureController does not implement SecuredCrudController (missing method GetAuthChecker)
```

#### Example 2: Using the Builder Without Required Methods

```go
// Try to build without auth:
controller := NewInterfaceEnforcedBuilder[Model, *CreateReq, *UpdateReq]("item", service).
    ValidateCreateRequest().
    ValidateUpdateRequest().
    // Forgot WithAuthChecker()!
    Build()

// RUNTIME PANIC (caught during development):
// panic: InterfaceEnforcedBuilder: WithAuthChecker() must be called before Build()
```

#### Example 3: Returning nil from Interface Methods

```go
type BadController struct {
    *BaseController
}

func (c *BadController) GetAuthChecker() func(...) error {
    return nil // Oops!
}

// This compiles, but when registering:
registrar.RegisterSecuredCRUD("/api/items", controller, controller)

// RUNTIME PANIC (caught during development):
// panic: Controller auth checker returned nil
```

## Benefits

1. **Compile-Time Safety**: You cannot register insecure controllers
2. **Clear Contracts**: Interfaces document security requirements
3. **IDE Support**: Autocomplete shows what methods are needed
4. **Refactoring Safety**: Changing interfaces breaks non-compliant code

## Migration Guide

To migrate existing controllers:

```go
// Before:
type BookController struct {
    *contracts.EnforcedCrudController[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest]
}

// After:
type BookController struct {
    *contracts.InterfaceEnforcedController[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest]
}

// Build with required methods:
controller := contracts.NewInterfaceEnforcedBuilder[...](...).
    ValidateCreateRequest().    // Required
    ValidateUpdateRequest().    // Required
    WithAuthChecker(...).       // Required
    Build()

// Register with type safety:
registrar := contracts.NewSecureRouteRegistrar(route)
registrar.RegisterSecuredCRUD("/api/books", controller, controller)
```

## Testing

Always verify your controllers implement the interfaces:

```go
// In your test file:
func TestBookController_ImplementsSecurityInterfaces(t *testing.T) {
    controller := NewBookController(mockService)
    
    // Compile-time check
    var _ contracts.SecuredCrudController = controller
    
    // Runtime checks
    assert.NotNil(t, controller.GetAuthChecker())
    assert.NotNil(t, controller.GetCreateValidator())
    assert.NotNil(t, controller.GetUpdateValidator())
}
```

## Summary

This pattern makes it **impossible** to accidentally create insecure controllers. The compiler becomes your security guard, refusing to build code that doesn't meet security requirements.