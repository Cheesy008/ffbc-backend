package service

import (
	"context"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

func (h *Handler) DetailCreation(
	ctx context.Context,
	input *DetailInput,
) (*CreationOutput, error) {
	result, err := h.useCase.GetByID(ctx, input.ID, catalogdomain.ServiceTypeCreation)
	if err != nil {
		return nil, mapError(err)
	}
	return &CreationOutput{Body: mapCreationResponse(result)}, nil
}

func (h *Handler) DetailSelling(
	ctx context.Context,
	input *DetailInput,
) (*SellingOutput, error) {
	result, err := h.useCase.GetByID(ctx, input.ID, catalogdomain.ServiceTypeSelling)
	if err != nil {
		return nil, mapError(err)
	}
	return &SellingOutput{Body: mapSellingResponse(result)}, nil
}
