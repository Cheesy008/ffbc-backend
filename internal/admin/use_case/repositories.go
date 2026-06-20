package use_case

import (
	"context"
	"time"

	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
)

type AdminUserRepository interface {
	Create(ctx context.Context, input AdminUserCreateInput) (domain.AdminUser, error)
	FindByEmail(ctx context.Context, email string) (domain.AdminUser, error)
	FindById(ctx context.Context, id int64) (domain.AdminUser, error)
}

type AdminSessionRepository interface {
	Create(ctx context.Context, session AdminSessionCreateInput) error
	FindByTokenHash(ctx context.Context, tokenHash string) (domain.AdminSession, error)
	Revoke(ctx context.Context, session AdminSessionRevokeInput) error
}

type AdminSessionCreateInput struct {
	AdminUserID int64
	TokenHash   string
	ExpiresAt   time.Time
}

type AdminSessionRevokeInput struct {
	TokenHash string
	RevokedAt time.Time
}

type AdminUserCreateInput struct {
	Email       string
	Password    string
	DisplayName *string
	IsActive    bool
}
