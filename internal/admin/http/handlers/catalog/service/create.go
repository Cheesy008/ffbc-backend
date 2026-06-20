package service

import (
	"context"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

func (h *Handler) CreateSelling(
	ctx context.Context,
	input *SellingCreateInput,
) (*SellingOutput, error) {
	serviceInput, err := sellingCreateInput(input.Body)
	if err != nil {
		return nil, err
	}

	result, err := h.useCase.Create(ctx, serviceInput)
	if err != nil {
		return nil, mapError(err)
	}

	return &SellingOutput{Body: mapSellingResponse(result)}, nil
}

func (h *Handler) CreateCreation(
	ctx context.Context,
	input *CreationCreateInput,
) (*CreationOutput, error) {
	serviceInput, err := creationCreateInput(input.Body)
	if err != nil {
		return nil, err
	}

	result, err := h.useCase.Create(ctx, serviceInput)
	if err != nil {
		return nil, mapError(err)
	}
	if result.Type != catalogdomain.ServiceTypeCreation {
		return nil, mapError(catalogdomain.ErrServiceTypeMismatch)
	}

	return &CreationOutput{Body: mapCreationResponse(result)}, nil
}
