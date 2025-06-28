package contracts

// ListRequest for pagination, sorting, and filtering
type ListRequest struct {
	Page      int                    `form:"page" json:"page"`
	PageSize  int                    `form:"pageSize" json:"pageSize"`
	Sort      string                 `form:"sort" json:"sort"`
	Direction string                 `form:"direction" json:"direction"`
	Search    string                 `form:"search" json:"search"`
	Filters   map[string]interface{} `form:"filters" json:"filters"`
}

// ListResponse for paginated results
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

// PaginatedResult represents paginated data
type PaginatedResult struct {
	Data        []interface{} `json:"data"`
	Total       int64         `json:"total"`
	PerPage     int           `json:"perPage"`
	CurrentPage int           `json:"currentPage"`
	LastPage    int           `json:"lastPage"`
	From        int           `json:"from"`
	To          int           `json:"to"`
	HasNext     bool          `json:"hasNext"`
	HasPrev     bool          `json:"hasPrev"`
}

// SetDefaults applies sensible defaults to ListRequest
func (r *ListRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 20
	}
	if r.PageSize > 100 {
		r.PageSize = 100
	}
	if r.Sort == "" {
		r.Sort = "id"
	}
	if r.Direction == "" {
		r.Direction = "DESC"
	}
	// Additional validation to prevent issues
	if r.PageSize == 0 {
		r.PageSize = 20
	}
}