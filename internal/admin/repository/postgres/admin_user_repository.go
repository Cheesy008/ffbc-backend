package admin_postgres

import (
	"context"
	"errors"

	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
	adminsqlc "github.com/cheesy008/ffbc-backend/internal/admin/repository/postgres/sqlc/generated"
	"github.com/cheesy008/ffbc-backend/internal/admin/use_case"
	"github.com/jackc/pgx/v5/pgconn"
)

type AdminUserRepository struct {
	queries *adminsqlc.Queries
}

func NewAdminUserRepository(queries *adminsqlc.Queries) *AdminUserRepository {
	return &AdminUserRepository{queries: queries}
}

func (r *AdminUserRepository) Create(ctx context.Context, input use_case.AdminUserCreateInput) (domain.AdminUser, error) {
	user, err := r.queries.CreateAdminUser(ctx, adminsqlc.CreateAdminUserParams{
		Email:        input.Email,
		PasswordHash: input.Password,
		DisplayName:  input.DisplayName,
		IsActive:     input.IsActive,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return domain.AdminUser{}, domain.ErrAdminAlreadyExists
		}

		return domain.AdminUser{}, err
	}

	return mapAdminUser(user), nil
}

func (r *AdminUserRepository) FindByEmail(ctx context.Context, email string) (domain.AdminUser, error) {
	user, err := r.queries.GetAdminUserByEmail(ctx, email)
	if err != nil {
		return domain.AdminUser{}, domain.ErrAdminNotFound
	}

	return mapAdminUser(user), nil
}

func (r *AdminUserRepository) FindById(ctx context.Context, id int64) (domain.AdminUser, error) {
	user, err := r.queries.GetAdminUserById(ctx, id)
	if err != nil {
		return domain.AdminUser{}, domain.ErrAdminNotFound
	}

	return mapAdminUser(user), nil
}

func mapAdminUser(user adminsqlc.AdminUser) domain.AdminUser {
	return domain.AdminUser{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		DisplayName:  user.DisplayName,
		IsActive:     user.IsActive,
		LastLoginAt:  user.LastLoginAt,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
