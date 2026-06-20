package use_case

import (
	"context"
	"fmt"
	"time"

	"github.com/cheesy008/ffbc-backend/internal/security/session_token"
)

type LogoutInput struct {
	SessionToken string
}

func (uc *AuthUseCase) Logout(ctx context.Context, input LogoutInput) error {
	tokenHash := session_token.Hash(input.SessionToken)
	fmt.Println(tokenHash)
	err := uc.sessionRepo.Revoke(ctx, AdminSessionRevokeInput{
		TokenHash: tokenHash,
		RevokedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}
