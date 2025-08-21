package tests

import (
	"testing"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

func TestStaticControllerCompile(t *testing.T) {
	service := services.NewBookService()
	
	// TEST 1: This should NOT compile - missing all methods
	contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		service,
	).Build() // ERROR: Build undefined (type *contracts.StaticControllerBuilder[...needsCreateRequest] has no field or method Build)
	
	// TEST 2: This should NOT compile - missing ValidateUpdateRequest and WithAuthChecker
	contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		service,
	).
		ValidateCreateRequest().
		Build() // ERROR: Build undefined (type *contracts.StaticControllerBuilder[...needsUpdateRequest] has no field or method Build)
	
	// TEST 3: This should NOT compile - missing WithAuthChecker
	contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		service,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
		Build() // ERROR: Build undefined (type *contracts.StaticControllerBuilder[...needsAuthChecker] has no field or method Build)
}