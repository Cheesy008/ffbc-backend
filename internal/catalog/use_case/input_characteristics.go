package use_case

import (
	"context"

	"github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	"github.com/cheesy008/ffbc-backend/internal/shared/optional"
)

type InputCharacteristicsUseCase struct {
	repo InputCharacteristicsRepository
}

func NewInputCharacteristicsUseCase(repo InputCharacteristicsRepository) *InputCharacteristicsUseCase {
	return &InputCharacteristicsUseCase{repo: repo}
}

type InputCharacteristicsInput struct {
	Name string
	Type domain.InputCharacteristicType
}

type InputCharacteristicsPatchInput struct {
	ID             int
	Name           optional.Optional[string]
	Type           optional.Optional[domain.InputCharacteristicType]
	TemplateIDList optional.Optional[*[]int]
}

type InputCharacteristicsListFilter struct {
	Search string
	Offset int
	Count  int
}

func (uc *InputCharacteristicsUseCase) BulkCreate(
	ctx context.Context,
	inputs []InputCharacteristicsInput,
) ([]domain.InputCharacteristic, error) {
	characteristics, err := uc.repo.BulkCreate(ctx, inputs)
	if err != nil {
		return nil, err
	}
	return characteristics, nil
}

func (uc *InputCharacteristicsUseCase) Patch(ctx context.Context, input InputCharacteristicsPatchInput) (domain.InputCharacteristic, error) {
	if !input.Name.Set && !input.Type.Set && !input.TemplateIDList.Set {
		return domain.InputCharacteristic{}, domain.ErrEmptyPatch
	}

	characteristic, err := uc.repo.Patch(ctx, input)
	if err != nil {
		return domain.InputCharacteristic{}, err
	}
	return characteristic, nil
}

func (uc *InputCharacteristicsUseCase) GetByID(ctx context.Context, id int) (domain.InputCharacteristic, error) {
	characteristic, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domain.InputCharacteristic{}, err
	}
	return characteristic, nil
}

func (uc *InputCharacteristicsUseCase) List(
	ctx context.Context,
	filter InputCharacteristicsListFilter,
) ([]domain.InputCharacteristic, error) {
	characteristics, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return characteristics, nil
}
