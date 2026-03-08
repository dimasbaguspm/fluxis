package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/user/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrEmailTaken   = httpx.Conflict("email already registerd")
	ErrUserNotFound = httpx.NotFound("user not found")
)

func (s *Service) GetSingleUserById(ctx context.Context, id pgtype.UUID) (domain.UserModel, error) {
	user, err := s.Repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, ErrUserNotFound
		}
		return domain.UserModel{}, fmt.Errorf("get user by id: %w", err)
	}

	return domain.UserModel{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Password:    user.PasswordHash,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
	}, nil
}

func (s *Service) GetSingleUserByEmail(ctx context.Context, email string) (domain.UserModel, error) {
	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, ErrUserNotFound
		}
		return domain.UserModel{}, fmt.Errorf("get user by email: %w", err)
	}
	return domain.UserModel{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Password:    user.PasswordHash,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, p domain.UserCreateModel) (domain.UserModel, error) {
	user, err := s.Repo.CreateUser(ctx, repository.CreateUserParams{
		Email:       p.Email,
		DisplayName: p.DisplayName,
		// hash handled by auth service
		PasswordHash: p.Password,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.UserModel{}, ErrEmailTaken
		}

		return domain.UserModel{}, fmt.Errorf("create user: %w", err)
	}

	return domain.UserModel{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Password:    user.PasswordHash,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, id pgtype.UUID, p domain.UserUpdateModel) (domain.UserModel, error) {
	user, err := s.Repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:      id,
		Column1: p.DisplayName,
		// hash handled by auth service
		Column2: p.Password,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, ErrUserNotFound
		}
		return domain.UserModel{}, fmt.Errorf("update user: %w", err)
	}

	return domain.UserModel{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Password:    user.PasswordHash,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
	}, nil

}

func (s *Service) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	err := s.Repo.DeleteUser(ctx, id)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
