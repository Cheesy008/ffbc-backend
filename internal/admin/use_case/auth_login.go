package use_case

import (
	"context"
	"time"

	"github.com/cheesy008/ffbc-backend/internal/admin/constants"
	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
	"github.com/cheesy008/ffbc-backend/internal/security/password"
	"github.com/cheesy008/ffbc-backend/internal/security/session_token"
)

type LoginInput struct {
	Email         string
	PlainPassword string
}

type LoginOutput struct {
	SessionToken string
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	user, err := uc.adminRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return LoginOutput{}, domain.ErrAdminNotFound
	}

	if !user.CanLogin() {
		return LoginOutput{}, domain.ErrAdminInactive
	}

	err = password.Compare(user.PasswordHash, input.PlainPassword)
	if err != nil {
		return LoginOutput{}, domain.ErrInvalidCredentials
	}

	plainToken, tokenHash, err := session_token.New()
	if err != nil {
		return LoginOutput{}, err
	}

	err = uc.sessionRepo.Create(ctx, AdminSessionCreateInput{
		AdminUserID: user.ID,
		TokenHash:   tokenHash,
		ExpiresAt:   time.Now().Add(constants.CookieExpirationTime),
	})
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{SessionToken: plainToken}, nil
}
