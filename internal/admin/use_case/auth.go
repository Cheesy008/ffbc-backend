package use_case

type AuthUseCase struct {
	adminRepo   AdminUserRepository
	sessionRepo AdminSessionRepository
}

func NewAuthUseCase(
	adminRepo AdminUserRepository,
	sessionRepo AdminSessionRepository) *AuthUseCase {
	return &AuthUseCase{
		adminRepo:   adminRepo,
		sessionRepo: sessionRepo,
	}
}
