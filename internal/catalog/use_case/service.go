package use_case

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	"github.com/cheesy008/ffbc-backend/internal/shared/optional"
)

type ServiceUseCase struct {
	repo ServiceRepository
}

type ServiceInputCharacteristicInput struct {
	InputCharacteristicID int
	IsRequired            bool
	SortOrder             *int
}

type ServiceInput struct {
	Name                 string
	BasePrice            domain.Decimal
	Description          *string
	Type                 domain.ServiceType
	Status               domain.ServiceStatus
	InputCharacteristics []ServiceInputCharacteristicInput
}

type ServiceListFilter struct {
	Search string
	Offset int
	Count  int
	Type   domain.ServiceType
}

type ServicePatchInput struct {
	ID                     int
	Type                   domain.ServiceType
	Name                   optional.Optional[string]
	BasePrice              optional.Optional[domain.Decimal]
	Description            optional.Optional[*string]
	InputCharacteristicIDs optional.Optional[[]int]
}

func NewServiceUseCase(repo ServiceRepository) *ServiceUseCase {
	return &ServiceUseCase{repo: repo}
}

func (uc *ServiceUseCase) Create(ctx context.Context, input ServiceInput) (domain.Service, error) {
	return uc.repo.Create(ctx, input)
}

func (uc *ServiceUseCase) Patch(ctx context.Context, input ServicePatchInput) (domain.Service, error) {
	if !input.Name.Set &&
		!input.BasePrice.Set &&
		!input.Description.Set &&
		!input.InputCharacteristicIDs.Set {
		return domain.Service{}, domain.ErrEmptyPatch
	}

	return uc.repo.Patch(ctx, input)
}

func (uc *ServiceUseCase) GetByID(
	ctx context.Context,
	id int,
	serviceType domain.ServiceType,
) (domain.Service, error) {
	return uc.repo.GetByID(ctx, id, serviceType)
}

func (uc *ServiceUseCase) List(
	ctx context.Context,
	filter ServiceListFilter,
) ([]domain.ServiceListItem, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *ServiceUseCase) Delete(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
