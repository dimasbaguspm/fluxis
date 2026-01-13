package models

import "time"

type StatusModel struct {
	ID        string     `json:"id" format:"uuid"`
	ProjectID string     `json:"projectId" format:"uuid"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Position  int        `json:"position"`
	IsDefault bool       `json:"isDefault"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type StatusCreateModel struct {
	ProjectID string `json:"projectId" format:"uuid" required:"true"`
	Name      string `json:"name" minLength:"1"`
}

type StatusUpdateModel struct {
	Name string `json:"name" minLength:"1"`
}

type StatusReorderModel struct {
	ProjectID string   `json:"projectId" format:"uuid" required:"true"`
	IDs       []string `json:"ids" format:"uuid"`
}
