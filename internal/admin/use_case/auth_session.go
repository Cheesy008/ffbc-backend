package use_case

import (
	"context"
	"time"

	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
	"github.com/cheesy008/ffbc-backend/internal/security/session_token"
)

func (uc *AuthUseCase) GetAdminBySessionToken(ctx context.Context, sessionToken string) (domain.AdminUser, error) {
	session, err := uc.sessionRepo.FindByTokenHash(ctx, session_token.Hash(sessionToken))
	if err != nil {
		return domain.AdminUser{}, domain.ErrorSessionNotFound
	}

	if !session.IsActive(time.Now()) {
		return domain.AdminUser{}, domain.ErrorSessionExpired
	}

	user, err := uc.adminRepo.FindById(ctx, session.AdminUserID)
	if err != nil {
		return domain.AdminUser{}, err
	}

	return user, nil
}
