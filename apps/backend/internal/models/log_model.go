package models

import "time"

type LogModel struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	TaskID    *string   `json:"taskId,omitempty"`
	StatusID  *string   `json:"statusId,omitempty"`
	Entry     string    `json:"entry"`
	CreatedAt time.Time `json:"createdAt"`
}

type LogCreateModel struct {
	ProjectID string  `json:"projectId" minLength:"1"`
	TaskID    *string `json:"taskId,omitempty"`
	StatusID  *string `json:"statusId,omitempty"`
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
	TaskID     []string `query:"taskId"`
	StatusID   []string `query:"statusId"`
	Query      string   `query:"query"`
	PageNumber int      `query:"pageNumber" default:"1"`
	PageSize   int      `query:"pageSize" default:"25"`
}
