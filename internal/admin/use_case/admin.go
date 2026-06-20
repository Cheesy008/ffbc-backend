package use_case

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
)

type AdminUseCase struct {
	adminRepo AdminUserRepository
}

func NewAdminUseCase(adminRepo AdminUserRepository) *AdminUseCase {
	return &AdminUseCase{adminRepo: adminRepo}
}

func (uc *AdminUseCase) Create(ctx context.Context, input AdminUserCreateInput) (domain.AdminUser, error) {
	user, err := uc.adminRepo.Create(ctx, input)
	if err != nil {
		return domain.AdminUser{}, err
	}

	return user, nil
}
