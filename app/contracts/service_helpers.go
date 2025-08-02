package contracts

import (
	"fmt"
	"github.com/goravel/framework/facades"
)

// SetActualServiceHelper is a helper function to set the actual service instance
// on a CrudServiceContract that supports it. This encapsulates the type assertion
// and error logging that would otherwise be repeated in every service constructor.
//
// This is necessary when using the ServiceBuilder pattern because the builder returns
// a generic CrudServiceContract, but the service needs a reference to the actual
// service implementation for proper method resolution when using embedded structs.
//
// Usage:
//   service := contracts.NewServiceBuilder[models.Book]("book", "id").Build()
//   bookServiceInstance := &BookService{
//       CrudServiceContract: service,
//       baseService: service,
//   }
//   contracts.SetActualServiceHelper(service, bookServiceInstance, "BookService")
func SetActualServiceHelper(service CrudServiceContract, actualService interface{}, serviceName string) {
	if setter, ok := service.(interface {
		SetActualService(interface{})
	}); ok {
		setter.SetActualService(actualService)
	} else {
		facades.Log().Error(fmt.Sprintf("%s: Failed to cast service to SetActualService interface", serviceName), map[string]interface{}{
			"serviceType": fmt.Sprintf("%T", service),
		})
	}
}