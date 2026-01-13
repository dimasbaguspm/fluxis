package models

import "time"

type TaskModel struct {
	ID        string     `json:"id" format:"uuid"`
	ProjectID string     `json:"projectId" format:"uuid"`
	StatusID  string     `json:"statusId" format:"uuid"`
	Title     string     `json:"title"`
	Details   string     `json:"details"`
	Priority  int        `json:"priority"`
	DueDate   *time.Time `json:"dueDate,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type TaskPaginatedModel struct {
	Items      []TaskModel `json:"items"`
	PageNumber int         `json:"pageNumber"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
	TotalCount int         `json:"totalCount"`
}

type TaskSearchModel struct {
	ID         []string `query:"id" format:"uuid"`
	ProjectID  []string `query:"projectId" format:"uuid"`
	StatusID   []string `query:"statusId" format:"uuid"`
	Query      string   `query:"query"`
	PageNumber int      `query:"pageNumber" default:"1"`
	PageSize   int      `query:"pageSize" default:"25"`
	SortBy     string   `query:"sortBy" enum:"dueDate,createdAt,updatedAt,priority" default:"dueDate"`
	SortOrder  string   `query:"sortOrder" enum:"asc,desc" default:"desc"`
}

type TaskCreateModel struct {
	ProjectID string     `json:"projectId" minLength:"1" format:"uuid"`
	StatusID  string     `json:"statusId" required:"true" format:"uuid"`
	Title     string     `json:"title" minLength:"1" pattern:"^.*\\S.*$"`
	Details   string     `json:"details"`
	Priority  int        `json:"priority" default:"1" minimum:"1"`
	DueDate   *time.Time `json:"dueDate,omitempty"`
}

type TaskUpdateModel struct {
	Title    string     `json:"title,omitempty" required:"false" minLength:"1" pattern:"^.*\\S.*$"`
	Details  string     `json:"details,omitempty" required:"false"`
	StatusID string     `json:"statusId,omitempty" required:"false"`
	Priority *int       `json:"priority,omitempty" required:"false" minimum:"1"`
	DueDate  *time.Time `json:"dueDate,omitempty" required:"false"`
}
