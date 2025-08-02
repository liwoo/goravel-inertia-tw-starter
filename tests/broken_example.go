package tests

import (
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// This file intentionally has a compile error to demonstrate enforcement

func BrokenExample() {
	service := services.NewBookService()
	
	// THIS SHOULD NOT COMPILE - missing WithAuthChecker
	_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		service,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
		// WithAuthChecker(...) is missing!
		Build() // This line should cause error: Build undefined
}