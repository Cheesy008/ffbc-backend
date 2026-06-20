package service

import (
	"context"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

func (h *Handler) PatchCreation(
	ctx context.Context,
	input *CreationPatchInput,
) (*CreationOutput, error) {
	serviceInput, err := creationPatchInput(input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	result, err := h.useCase.Patch(ctx, serviceInput)
	if err != nil {
		return nil, mapError(err)
	}
	return &CreationOutput{Body: mapCreationResponse(result)}, nil
}

func (h *Handler) PatchSelling(
	ctx context.Context,
	input *SellingPatchInput,
) (*SellingOutput, error) {
	serviceInput, err := sellingPatchInput(input.ID, input.Body)
	if err != nil {
		return nil, err
	}
	serviceInput.Type = catalogdomain.ServiceTypeSelling

	result, err := h.useCase.Patch(ctx, serviceInput)
	if err != nil {
		return nil, mapError(err)
	}
	return &SellingOutput{Body: mapSellingResponse(result)}, nil
}
