package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID          string    `json:"id" validate:"required,uuid4"`
	Email       string    `json:"email" validate:"email"`
	DisplayName string    `json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserPagedModel struct {
	Items      []UserModel `json:"-"`
	PageSize   int         `json:"pageSize"`
	PageItems  int         `json:"pageItems"`
	TotalPages int         `json:"totalPages"`
	TotalItems int         `json:"totalItems"`
}

type UserSearchModel struct {
	IDs         string `json:"ids" validate:"dive,uuid4"`
	Email       string `json:"email" validate:"email"`
	DisplayName string `json:"displayName"`
}

type UserCreateModel struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateModel struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password" validate:"required"`
}

type UserRead interface {
	GetPagedUsers(ctx context.Context, q UserSearchModel) (UserPagedModel, error)
	GetSingleUser(ctx context.Context, id uuid.UUID) (UserModel, error)
}

type UserWrite interface {
	CreateUser(ctx context.Context, p UserCreateModel) (UserModel, error)
	UpdateUser(ctx context.Context, p UserUpdateModel) (UserModel, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
