package inputcharacteristics

import (
	"context"

	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
	"github.com/cheesy008/ffbc-backend/internal/shared/listquery"
)

type ListInput struct {
	Search string `query:"search" required:"false" example:"Длина"`
	Offset int    `query:"offset" minimum:"0" default:"0" example:"0"`
	Count  int    `query:"count" minimum:"1" maximum:"100" default:"20" example:"20"`
}

type ListOutput struct {
	Body []InputCharacteristicListItemResponse
}

func (h *Handler) List(
	ctx context.Context,
	input *ListInput,
) (*ListOutput, error) {
	filter, err := listquery.Normalize(listquery.Input{
		Search: input.Search,
		Offset: input.Offset,
		Count:  input.Count,
	}, listquery.Options{})
	if err != nil {
		return nil, mapListQueryError(err)
	}

	characteristics, err := h.useCase.List(ctx, catalogusecase.InputCharacteristicsListFilter{
		Search: filter.Search,
		Offset: filter.Offset,
		Count:  filter.Count,
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]InputCharacteristicListItemResponse, 0, len(characteristics))
	for _, characteristic := range characteristics {
		items = append(items, mapListItemResponse(characteristic))
	}

	return &ListOutput{Body: items}, nil
}
