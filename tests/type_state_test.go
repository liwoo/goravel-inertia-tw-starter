package tests

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
	"testing"
)

func TestTypeStatePattern(t *testing.T) {
	service := services.NewBookService()

	// Get the builder at different stages to inspect types
	builder1 := contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		service,
	)
	// Type of builder1: *StaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest, needsCreateRequest]

	builder2 := builder1.ValidateCreateRequest()
	// Type of builder2: *StaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest, needsUpdateRequest]

	builder3 := builder2.ValidateUpdateRequest()
	// Type of builder3: *StaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest, needsAuthChecker]

	builder4 := builder3.WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
		return nil
	})
	// Type of builder4: *StaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest, readyToBuild]

	controller := builder4.Build()
	// Only builder4 (readyToBuild state) should have Build() method

	_ = controller

	// These should NOT compile if uncommented:
	// builder1.Build() // Error: Build undefined
	// builder2.Build() // Error: Build undefined
	// builder3.Build() // Error: Build undefined

	// builder1.ValidateUpdateRequest() // Error: ValidateUpdateRequest undefined
	// builder1.WithAuthChecker(...)    // Error: WithAuthChecker undefined

	// builder2.WithAuthChecker(...) // Error: WithAuthChecker undefined
}
