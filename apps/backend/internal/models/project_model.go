package models

import "time"

type ProjectModel struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status" enum:"active,paused,archived"`
	CreatedAt   time.Time  `json:"createdAt" enum:"active,paused,archived"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type ProjectPaginatedModel struct {
	Items      []ProjectModel `json:"items"`
	PageNumber int            `json:"pageNumber"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
	TotalCount int            `json:"totalCount"`
}

type ProjectSearchModel struct {
	ID         []string `query:"id"`
	Query      string   `query:"query"`
	Status     []string `query:"status" enum:"active,paused,archived"`
	PageNumber int      `query:"pageNumber" default:"1"`
	PageSize   int      `query:"pageSize" default:"25"`
	SortBy     string   `query:"sortBy" enum:"createdAt,updatedAt,status" default:"createdAt"`
	SortOrder  string   `query:"sortOrder" enum:"asc,desc" default:"desc"`
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
