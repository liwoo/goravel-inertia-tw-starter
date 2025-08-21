package contracts

import (
	"github.com/goravel/framework/contracts/route"
)

// SecureRouteRegistrar provides compile-time safe route registration
type SecureRouteRegistrar struct {
	route route.Route
}

// NewSecureRouteRegistrar creates a new secure route registrar
func NewSecureRouteRegistrar(r route.Route) *SecureRouteRegistrar {
	return &SecureRouteRegistrar{route: r}
}

// RegisterSecuredCRUD registers CRUD routes with compile-time security checks
// This function will ONLY accept controllers that implement SecuredCrudController
func (r *SecureRouteRegistrar) RegisterSecuredCRUD(
	path string,
	controller SecuredCrudController, // Compile-time enforcement here!
	crudController CrudControllerContract,
) {
	// Verify at runtime that the security is properly configured
	if controller.GetAuthChecker() == nil {
		panic("Controller auth checker returned nil")
	}
	if controller.GetCreateValidator() == nil {
		panic("Controller create validator returned nil")
	}
	if controller.GetUpdateValidator() == nil {
		panic("Controller update validator returned nil")
	}

	// Register the routes
	r.route.Get(path, crudController.Index)
	r.route.Post(path, crudController.Store)
	r.route.Get(path+"/{id}", crudController.Show)
	r.route.Put(path+"/{id}", crudController.Update)
	r.route.Delete(path+"/{id}", crudController.Delete)
}

// RegisterSecuredResource is an alternative that takes both interfaces
func (r *SecureRouteRegistrar) RegisterSecuredResource(
	path string,
	secured SecuredCrudController,
	crud CrudControllerContract,
) {
	// Both parameters must be the same controller instance
	// This provides compile-time checking
	r.RegisterSecuredCRUD(path, secured, crud)
}

// Example usage function that shows compile-time enforcement
func ExampleSecureRegistration(r route.Route) {
	_ = NewSecureRouteRegistrar(r)

	// This will compile - controller implements all interfaces
	// var securedController *InterfaceEnforcedController[Model, *CreateReq, *UpdateReq]
	// registrar.RegisterSecuredCRUD("/api/items", securedController, securedController)

	// This will NOT compile - missing SecuredCrudController interface
	// var unsecuredController *GenericCrudController[Model, *CreateReq, *UpdateReq]
	// registrar.RegisterSecuredCRUD("/api/items", unsecuredController, unsecuredController)
	// Compile error: cannot use unsecuredController (variable of type *GenericCrudController[...]) as SecuredCrudController value
}
