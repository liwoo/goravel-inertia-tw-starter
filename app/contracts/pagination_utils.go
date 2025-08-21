package contracts

import (
	"github.com/goravel/framework/support/collect"
)

// PaginationBuilder provides utilities for manual pagination
type PaginationBuilder struct{}

// NewPaginationBuilder creates a new pagination builder instance
func NewPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{}
}

// PaginateSlice performs manual pagination on a slice of any type
func (pb *PaginationBuilder) PaginateSlice(items interface{}, page, pageSize int) *PaginatedResult {
	// Use reflection to handle any slice type
	switch v := items.(type) {
	case []interface{}:
		return pb.paginateInterfaceSlice(v, page, pageSize)
	default:
		// For typed slices, we'll use a generic approach
		return pb.paginateGenericSlice(items, page, pageSize)
	}
}

// paginateInterfaceSlice handles []interface{} slices
func (pb *PaginationBuilder) paginateInterfaceSlice(items []interface{}, page, pageSize int) *PaginatedResult {
	total := int64(len(items))
	offset := (page - 1) * pageSize
	end := offset + pageSize

	// Ensure offset is within bounds
	if offset > len(items) {
		offset = len(items)
	}
	if end > len(items) {
		end = len(items)
	}

	// Extract page items
	var pageItems []interface{}
	if offset < len(items) {
		pageItems = items[offset:end]
	} else {
		pageItems = []interface{}{}
	}

	// Calculate pagination metadata
	lastPage := pb.CalculateLastPage(total, int64(pageSize))

	return &PaginatedResult{
		Data:        pageItems,
		Total:       total,
		PerPage:     pageSize,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len(pageItems),
	}
}

// paginateGenericSlice handles any typed slice using a callback approach
func (pb *PaginationBuilder) paginateGenericSlice(items interface{}, page, pageSize int) *PaginatedResult {
	// This is a placeholder - in practice, services will use PaginateSliceWithConverter
	return &PaginatedResult{
		Data:        []interface{}{},
		Total:       0,
		PerPage:     pageSize,
		CurrentPage: page,
		LastPage:    1,
		From:        0,
		To:          0,
	}
}

// PaginateSliceWithConverter performs pagination and converts items to []interface{}
func PaginateSliceWithConverter[T any](items []T, page, pageSize int, converter func(T) interface{}) *PaginatedResult {
	total := int64(len(items))
	offset := (page - 1) * pageSize
	end := offset + pageSize

	// Ensure offset is within bounds
	if offset > len(items) {
		offset = len(items)
	}
	if end > len(items) {
		end = len(items)
	}

	// Extract page items
	var pageItems []T
	if offset < len(items) {
		pageItems = items[offset:end]
	} else {
		pageItems = []T{}
	}

	// Convert to interface slice
	data := make([]interface{}, len(pageItems))
	if converter != nil {
		for i, item := range pageItems {
			data[i] = converter(item)
		}
	} else {
		// Default converter - just cast to interface{}
		for i, item := range pageItems {
			data[i] = item
		}
	}

	// Calculate pagination metadata
	pb := &PaginationBuilder{}
	lastPage := pb.CalculateLastPage(total, int64(pageSize))

	return &PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     pageSize,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        pb.CalculateFrom(offset, len(pageItems)),
		To:          pb.CalculateTo(offset, len(pageItems)),
	}
}

// ManualPaginate performs manual pagination on a slice and returns offset/limit
func ManualPaginate[T any](items []T, page, pageSize int) (pageItems []T, total int64, from int, to int) {
	total = int64(len(items))
	offset := (page - 1) * pageSize
	end := offset + pageSize

	// Ensure offset is within bounds
	if offset > len(items) {
		offset = len(items)
	}
	if end > len(items) {
		end = len(items)
	}

	// Extract page items
	if offset < len(items) {
		pageItems = items[offset:end]
	} else {
		pageItems = []T{}
	}

	pb := &PaginationBuilder{}
	from = pb.CalculateFrom(offset, len(pageItems))
	to = pb.CalculateTo(offset, len(pageItems))

	return pageItems, total, from, to
}

// CalculateLastPage calculates the last page number
func (pb *PaginationBuilder) CalculateLastPage(total, pageSize int64) int {
	if pageSize <= 0 {
		return 1
	}
	return int((total + pageSize - 1) / pageSize)
}

// CalculateOffset calculates the offset for a given page
func (pb *PaginationBuilder) CalculateOffset(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

// CalculateFrom calculates the "from" index for pagination metadata (1-based)
func (pb *PaginationBuilder) CalculateFrom(offset, itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	return offset + 1
}

// CalculateTo calculates the "to" index for pagination metadata (1-based)
func (pb *PaginationBuilder) CalculateTo(offset, itemCount int) int {
	if itemCount == 0 {
		return 0
	}
	return offset + itemCount
}

// ValidatePaginationBounds ensures offset and limit are within valid bounds
func (pb *PaginationBuilder) ValidatePaginationBounds(offset, limit, totalItems int) (validOffset, validLimit int) {
	// Ensure offset is not negative
	if offset < 0 {
		offset = 0
	}

	// Ensure offset doesn't exceed total items
	if offset >= totalItems {
		offset = totalItems
	}

	// Ensure limit is positive
	if limit <= 0 {
		limit = 20 // Default page size
	}

	// Calculate valid end position
	end := offset + limit
	if end > totalItems {
		end = totalItems
	}

	// Adjust limit based on actual available items
	validLimit = end - offset
	if validLimit < 0 {
		validLimit = 0
	}

	return offset, validLimit
}

// CreateEmptyResult creates an empty pagination result
func (pb *PaginationBuilder) CreateEmptyResult(page, pageSize int) *PaginatedResult {
	return &PaginatedResult{
		Data:        []interface{}{},
		Total:       0,
		PerPage:     pageSize,
		CurrentPage: page,
		LastPage:    1,
		From:        0,
		To:          0,
	}
}

// HasNextPage checks if there's a next page
func (pb *PaginationBuilder) HasNextPage(currentPage, lastPage int) bool {
	return currentPage < lastPage
}

// HasPrevPage checks if there's a previous page
func (pb *PaginationBuilder) HasPrevPage(currentPage int) bool {
	return currentPage > 1
}

// GetPageRange returns a range of page numbers for pagination UI
func (pb *PaginationBuilder) GetPageRange(currentPage, lastPage, maxVisible int) []int {
	if maxVisible <= 0 {
		maxVisible = 10
	}

	// Calculate start and end of range
	halfVisible := maxVisible / 2
	start := currentPage - halfVisible
	end := currentPage + halfVisible

	// Adjust boundaries
	if start < 1 {
		start = 1
		end = collect.Min([]int{maxVisible, lastPage})
	}
	if end > lastPage {
		end = lastPage
		start = collect.Max([]int{1, end - maxVisible + 1})
	}

	// Build page range
	pages := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}

	return pages
}
