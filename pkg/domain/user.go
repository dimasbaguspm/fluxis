package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserModel struct {
	ID          pgtype.UUID `json:"id" validate:"required,uuid4"`
	Email       string      `json:"email" validate:"email"`
	Password    string      `json:"password"`
	DisplayName string      `json:"displayName"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type UserPagedModel struct {
	Items      []UserModel `json:"-"`
	PageSize   int         `json:"pageSize"`
	PageItems  int         `json:"pageItems"`
	TotalPages int         `json:"totalPages"`
	TotalItems int         `json:"totalItems"`
}

type UserSearchModel struct {
	IDs         []pgtype.UUID `json:"ids" validate:"dive,uuid4"`
	Email       string        `json:"email" validate:"email"`
	DisplayName string        `json:"displayName" validate:"min=1"`
}

type UserCreateModel struct {
	Email       string `json:"email" validate:"email,required"`
	DisplayName string `json:"displayName" validate:"min=1"`
	Password    string `json:"password" validate:"required"`
}

type UserUpdateModel struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password" validate:"required"`
}

type UserRead interface {
	GetSingleUserById(ctx context.Context, id pgtype.UUID) (UserModel, error)
	GetSingleUserByEmail(ctx context.Context, email string) (UserModel, error)
}

type UserWrite interface {
	CreateUser(ctx context.Context, p UserCreateModel) (UserModel, error)
	UpdateUser(ctx context.Context, id pgtype.UUID, p UserUpdateModel) (UserModel, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}
