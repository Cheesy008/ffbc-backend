package use_case

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

type ServiceCategoryUseCase struct {
	repo ServiceCategoryRepository
}

type ServiceCategoryRepository interface {
	Create(ctx context.Context, input ServiceCategoryInput) (domain.ServiceCategory, error)
	Update(ctx context.Context, input ServiceCategoryInput) (domain.ServiceCategory, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter ServiceCategoryListFilter) ([]domain.ServiceCategory, error)
}

func NewServiceCategoryUseCase(repo ServiceCategoryRepository) *ServiceCategoryUseCase {
	return &ServiceCategoryUseCase{repo: repo}
}

type ServiceCategoryInput struct {
	ID   int
	Name string
}

type ServiceCategoryListFilter struct {
	Search    string
	Offset    int
	Count     int
	SortBy    string
	SortOrder string
}

func (uc *ServiceCategoryUseCase) Create(ctx context.Context, input ServiceCategoryInput) (domain.ServiceCategory, error) {
	category, err := uc.repo.Create(ctx, input)
	if err != nil {
		return domain.ServiceCategory{}, err
	}
	return category, nil
}

func (uc *ServiceCategoryUseCase) Update(ctx context.Context, input ServiceCategoryInput) (domain.ServiceCategory, error) {
	category, err := uc.repo.Update(ctx, input)
	if err != nil {
		return domain.ServiceCategory{}, err
	}
	return category, nil
}

func (uc *ServiceCategoryUseCase) Delete(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ServiceCategoryUseCase) List(
	ctx context.Context,
	filter ServiceCategoryListFilter,
) ([]domain.ServiceCategory, error) {
	categories, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
