package inputcharacteristics

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
)

type BulkCreateInput struct {
	Body []InputCharacteristicCreateRequest `minItems:"1"`
}

type BulkCreateOutput struct {
	Body []InputCharacteristicBulkResponse
}

func (h *Handler) BulkCreate(
	ctx context.Context,
	input *BulkCreateInput,
) (*BulkCreateOutput, error) {
	if len(input.Body) == 0 {
		return nil, huma.Error400BadRequest("at least one input characteristic is required")
	}

	inputs := make([]catalogusecase.InputCharacteristicsInput, 0, len(input.Body))
	for _, request := range input.Body {
		characteristicInput, err := createInput(request)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, characteristicInput)
	}

	characteristics, err := h.useCase.BulkCreate(ctx, inputs)
	if err != nil {
		return nil, mapError(err)
	}

	response := make([]InputCharacteristicBulkResponse, 0, len(characteristics))
	for _, characteristic := range characteristics {
		response = append(response, mapBulkResponse(characteristic))
	}

	return &BulkCreateOutput{Body: response}, nil
}
