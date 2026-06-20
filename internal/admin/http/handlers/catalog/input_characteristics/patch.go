package inputcharacteristics

import "context"

type PatchInput struct {
	ID   int `path:"id" minimum:"1" example:"1"`
	Body InputCharacteristicPatchRequest
}

type Output struct {
	Body InputCharacteristicResponse
}

func (h *Handler) Patch(
	ctx context.Context,
	input *PatchInput,
) (*Output, error) {
	characteristicInput, err := patchInput(input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	characteristic, err := h.useCase.Patch(ctx, characteristicInput)
	if err != nil {
		return nil, mapError(err)
	}

	return &Output{Body: mapResponse(characteristic)}, nil
}
