package use_case

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	"github.com/cheesy008/ffbc-backend/internal/shared/optional"
)

type InputCharacteristicsTemplateUseCase struct {
	repo InputCharacteristicsTemplateRepository
}

func NewInputCharacteristicsTemplateUseCase(repo InputCharacteristicsTemplateRepository) *InputCharacteristicsTemplateUseCase {
	return &InputCharacteristicsTemplateUseCase{repo: repo}
}

type InputCharacteristicsTemplateInput struct {
	ID                     int
	Name                   string
	Description            *string
	InputCharacteristicIDs []int
}

type InputCharacteristicsTemplatePatchInput struct {
	ID                     int
	Name                   optional.Optional[string]
	InputCharacteristicIDs optional.Optional[[]int]
	Description            optional.Optional[*string]
}

type InputCharacteristicsTemplateListFilter struct {
	Search                 string
	InputCharacteristicIDs []int
	Offset                 int
	Count                  int
	SortBy                 string
	SortOrder              string
}

func (uc *InputCharacteristicsTemplateUseCase) Create(
	ctx context.Context,
	input InputCharacteristicsTemplateInput,
) (domain.InputCharacteristicTemplate, error) {
	template, err := uc.repo.Create(ctx, input)
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}
	return template, nil
}

func (uc *InputCharacteristicsTemplateUseCase) Patch(ctx context.Context, input InputCharacteristicsTemplatePatchInput) (domain.InputCharacteristicTemplate, error) {
	if !input.Name.Set && !input.InputCharacteristicIDs.Set && !input.Description.Set {
		return domain.InputCharacteristicTemplate{}, domain.ErrEmptyPatch
	}

	template, err := uc.repo.Patch(ctx, input)
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}
	return template, nil
}

func (uc *InputCharacteristicsTemplateUseCase) GetByID(ctx context.Context, id int) (domain.InputCharacteristicTemplate, error) {
	template, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domain.InputCharacteristicTemplate{}, err
	}
	return template, nil
}

func (uc *InputCharacteristicsTemplateUseCase) Delete(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *InputCharacteristicsTemplateUseCase) List(
	ctx context.Context,
	filter InputCharacteristicsTemplateListFilter,
) ([]domain.InputCharacteristicTemplate, error) {
	templates, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return templates, nil
}
