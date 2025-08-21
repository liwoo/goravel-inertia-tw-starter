package contracts

// PaginatedResponse represents a strongly typed paginated API response
type PaginatedResponse struct {
	Data       interface{}            `json:"data"`
	Pagination PaginationInfo         `json:"pagination"`
	Filters    FiltersMeta            `json:"filters"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

// PaginationInfo represents pagination metadata in responses
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	From        int   `json:"from"`
	To          int   `json:"to"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// FiltersMeta represents the current filters/search/sort state in responses
type FiltersMeta struct {
	Page      int                    `json:"page"`
	PageSize  int                    `json:"pageSize"`
	Search    string                 `json:"search"`
	Sort      string                 `json:"sort"`
	Direction string                 `json:"direction"`
	Filters   map[string]interface{} `json:"filters"`
}

// SingleResourceResponse represents a response for a single resource
type SingleResourceResponse struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// CollectionResponse represents a response for a collection without pagination
type CollectionResponse struct {
	Data  []interface{}          `json:"data"`
	Total int                    `json:"total"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

// SearchResponse represents a search results response
type SearchResponse struct {
	Results    interface{}    `json:"results"`
	Query      string         `json:"query"`
	Pagination PaginationInfo `json:"pagination"`
	Filters    FiltersMeta    `json:"filters"`
	Meta       SearchMeta     `json:"meta,omitempty"`
}

// SearchMeta represents search-specific metadata
type SearchMeta struct {
	SearchTime   float64  `json:"search_time_ms,omitempty"`
	SearchedIn   []string `json:"searched_in,omitempty"`
	TotalMatches int      `json:"total_matches,omitempty"`
	Highlighted  bool     `json:"highlighted,omitempty"`
}

// BulkOperationResponse represents a response for bulk operations
type BulkOperationResponse struct {
	Success      []uint                 `json:"success"`
	Failed       []BulkOperationError   `json:"failed"`
	Total        int                    `json:"total"`
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
}

// BulkOperationError represents an error for a specific item in bulk operation
type BulkOperationError struct {
	ID     uint   `json:"id"`
	Error  string `json:"error"`
	Reason string `json:"reason,omitempty"`
}

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Errors  map[string][]string    `json:"errors"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Helper methods for PaginatedResponse

// NewPaginatedResponse creates a new paginated response from a PaginatedResult and ListRequest
func NewPaginatedResponse(result *PaginatedResult, request *ListRequest) *PaginatedResponse {
	return &PaginatedResponse{
		Data: result.Data,
		Pagination: PaginationInfo{
			CurrentPage: result.CurrentPage,
			LastPage:    result.LastPage,
			PerPage:     result.PerPage,
			Total:       result.Total,
			From:        result.From,
			To:          result.To,
			HasNext:     result.HasNext,
			HasPrev:     result.HasPrev,
		},
		Filters: FiltersMeta{
			Page:      request.Page,
			PageSize:  request.PageSize,
			Search:    request.Search,
			Sort:      request.Sort,
			Direction: request.Direction,
			Filters:   request.Filters,
		},
	}
}

// ToMap converts PaginatedResponse to map for compatibility
func (r *PaginatedResponse) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"data":       r.Data,
		"pagination": r.Pagination,
		"filters":    r.Filters,
	}

	if r.Meta != nil && len(r.Meta) > 0 {
		result["meta"] = r.Meta
	}

	return result
}

// WithMeta adds metadata to the response
func (r *PaginatedResponse) WithMeta(key string, value interface{}) *PaginatedResponse {
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	r.Meta[key] = value
	return r
}

// Helper methods for SearchResponse

// NewSearchResponse creates a new search response
func NewSearchResponse(results interface{}, query string, pagination PaginationInfo, filters FiltersMeta) *SearchResponse {
	return &SearchResponse{
		Results:    results,
		Query:      query,
		Pagination: pagination,
		Filters:    filters,
	}
}

// WithSearchMeta adds search metadata to the response
func (r *SearchResponse) WithSearchMeta(meta SearchMeta) *SearchResponse {
	r.Meta = meta
	return r
}

// Helper methods for BulkOperationResponse

// NewBulkOperationResponse creates a new bulk operation response
func NewBulkOperationResponse() *BulkOperationResponse {
	return &BulkOperationResponse{
		Success: []uint{},
		Failed:  []BulkOperationError{},
	}
}

// AddSuccess adds a successful operation
func (r *BulkOperationResponse) AddSuccess(id uint) {
	r.Success = append(r.Success, id)
	r.SuccessCount++
	r.Total++
}

// AddError adds a failed operation
func (r *BulkOperationResponse) AddError(id uint, err string, reason string) {
	r.Failed = append(r.Failed, BulkOperationError{
		ID:     id,
		Error:  err,
		Reason: reason,
	})
	r.FailedCount++
	r.Total++
}

// IsCompleteSuccess checks if all operations succeeded
func (r *BulkOperationResponse) IsCompleteSuccess() bool {
	return r.FailedCount == 0 && r.Total > 0
}

// Helper methods for ValidationErrorResponse

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(message string, errors map[string][]string) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

// AddError adds a validation error for a field
func (r *ValidationErrorResponse) AddError(field string, message string) {
	if r.Errors == nil {
		r.Errors = make(map[string][]string)
	}
	r.Errors[field] = append(r.Errors[field], message)
}

// HasErrors checks if there are any validation errors
func (r *ValidationErrorResponse) HasErrors() bool {
	return len(r.Errors) > 0
}
