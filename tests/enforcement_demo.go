package tests

import (
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
	"github.com/goravel/framework/contracts/http"
)

// Try commenting out any of the three methods below to see compile errors

func DemoCorrectUsage() {
	service := services.NewBookService()
	
	// This compiles correctly
	controller := contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		service,
	).
		ValidateCreateRequest().  // Try commenting this line
		ValidateUpdateRequest().  // Or try commenting this line
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			return nil
		}).                      // Or try commenting this line
		Build()
		
	_ = controller
}

// What happens when you comment out each line:
//
// 1. Comment out ValidateCreateRequest():
//    Next line fails: ValidateUpdateRequest undefined
//    
// 2. Comment out ValidateUpdateRequest():
//    Next line fails: WithAuthChecker undefined
//    
// 3. Comment out WithAuthChecker():
//    Next line fails: Build undefined
//
// The key is that each method returns a DIFFERENT TYPE:
// - NewStaticControllerBuilder returns StaticControllerBuilder[T,C,U,needsCreateRequest]
// - ValidateCreateRequest() returns StaticControllerBuilder[T,C,U,needsUpdateRequest]
// - ValidateUpdateRequest() returns StaticControllerBuilder[T,C,U,needsAuthChecker]
// - WithAuthChecker() returns StaticControllerBuilder[T,C,U,readyToBuild]
// - Only readyToBuild has the Build() method!