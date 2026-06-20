package admin_postgres

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
	adminsqlc "github.com/cheesy008/ffbc-backend/internal/admin/repository/postgres/sqlc/generated"
	"github.com/cheesy008/ffbc-backend/internal/admin/use_case"
)

type AdminSessionRepository struct {
	queries *adminsqlc.Queries
}

func NewAdminSessionRepository(queries *adminsqlc.Queries) *AdminSessionRepository {
	return &AdminSessionRepository{queries: queries}
}

func (r *AdminSessionRepository) Create(ctx context.Context, sessionIn use_case.AdminSessionCreateInput) error {
	_, err := r.queries.CreateAdminSession(ctx, adminsqlc.CreateAdminSessionParams{
		AdminUserID: sessionIn.AdminUserID,
		TokenHash:   sessionIn.TokenHash,
		ExpiresAt:   sessionIn.ExpiresAt,
	})
	if err != nil {
		return domain.ErrorSessionCreation
	}

	return nil
}

func (r *AdminSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (domain.AdminSession, error) {
	session, err := r.queries.GetAdminSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return domain.AdminSession{}, domain.ErrorSessionNotFound
	}

	return mapAdminSession(session), nil
}

func (r *AdminSessionRepository) Revoke(ctx context.Context, revokeIn use_case.AdminSessionRevokeInput) error {
	err := r.queries.RevokeAdminSession(ctx, adminsqlc.RevokeAdminSessionParams{
		RevokedAt: &revokeIn.RevokedAt,
		TokenHash: revokeIn.TokenHash,
	})
	if err != nil {
		return domain.ErrorSessionRevoke
	}

	return nil
}

func mapAdminSession(session adminsqlc.AdminSession) domain.AdminSession {
	return domain.AdminSession{
		ID:          session.ID,
		AdminUserID: session.AdminUserID,
		TokenHash:   session.TokenHash,
		ExpiresAt:   session.ExpiresAt,
		RevokedAt:   session.RevokedAt,
		CreatedAt:   session.CreatedAt,
		LastUsedAt:  session.LastUsedAt,
	}
}
