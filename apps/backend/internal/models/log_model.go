package models

import "time"

type LogModel struct {
	ID        string    `json:"id" format:"uuid"`
	ProjectID *string   `json:"projectId" format:"uuid"`
	TaskID    *string   `json:"taskId,omitempty" format:"uuid"`
	StatusID  *string   `json:"statusId,omitempty" format:"uuid"`
	Entry     string    `json:"entry"`
	CreatedAt time.Time `json:"createdAt"`
}

type LogCreateModel struct {
	ProjectID string  `json:"projectId" minLength:"1" format:"uuid"`
	TaskID    *string `json:"taskId,omitempty" format:"uuid"`
	StatusID  *string `json:"statusId,omitempty" format:"uuid"`
	Entry     string  `json:"entry" minLength:"1"`
}

type LogPaginatedModel struct {
	Items      []LogModel `json:"items"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	TotalPages int        `json:"totalPages"`
	TotalCount int        `json:"totalCount"`
}

type LogSearchModel struct {
	TaskID     []string `query:"taskId" format:"uuid"`
	StatusID   []string `query:"statusId" format:"uuid"`
	Query      string   `query:"query"`
	PageNumber int      `query:"pageNumber" default:"1"`
	PageSize   int      `query:"pageSize" default:"25"`
}
