package inputcharacteristics

import "context"

type DetailInput struct {
	ID int `path:"id" minimum:"1" example:"1"`
}

func (h *Handler) Detail(
	ctx context.Context,
	input *DetailInput,
) (*Output, error) {
	characteristic, err := h.useCase.GetByID(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}

	return &Output{Body: mapResponse(characteristic)}, nil
}
