package models

type ProjectModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	DeletedAt   string `json:"deletedAt,omitempty"`
}

type ProjectPaginatedModel struct {
	Items      []ProjectModel `json:"items"`
	PageNumber int            `json:"pageNumber"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
	TotalCount int            `json:"totalCount"`
	SortBy     string         `json:"sortBy" enum:"createdAt,updatedAt,status"`
	SortOrder  string         `json:"sortOrder" enum:"asc,desc"`
}

type ProjectSearchModel struct {
	ID         []string `json:"id"`
	Query      string   `json:"query"`
	Status     []string `json:"status" enum:"active,paused,archived"`
	PageNumber int      `json:"pageNumber" default:"1"`
	PageSize   int      `json:"pageSize" default:"25"`
	SortBy     string   `json:"sortBy" enum:"createdAt,updatedAt,status" default:"createdAt"`
	SortOrder  string   `json:"sortOrder" enum:"asc,desc" default:"desc"`
}

type ProjectCreateModel struct {
	Name        string `json:"name" minLength:"1"`
	Description string `json:"description" minLength:"1"`
	Status      string `json:"status" enum:"active,paused,archived"`
}

type ProjectUpdateModel struct {
	Name        string `json:"name,omitempty" required:"false" minLength:"1"`
	Description string `json:"description,omitempty" required:"false" minLenght:"1"`
	Status      string `json:"status,omitempty" required:"false" enum:"active,paused,archived"`
}
