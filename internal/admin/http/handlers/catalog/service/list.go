package service

import (
	"context"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
	"github.com/cheesy008/ffbc-backend/internal/shared/listquery"
)

func (h *Handler) ListCreation(ctx context.Context, input *ListInput) (*ListOutput, error) {
	return h.list(ctx, input, catalogdomain.ServiceTypeCreation)
}

func (h *Handler) ListSelling(ctx context.Context, input *ListInput) (*ListOutput, error) {
	return h.list(ctx, input, catalogdomain.ServiceTypeSelling)
}

func (h *Handler) list(
	ctx context.Context,
	input *ListInput,
	serviceType catalogdomain.ServiceType,
) (*ListOutput, error) {
	filter, err := listquery.Normalize(listquery.Input{
		Search: input.Search,
		Offset: input.Offset,
		Count:  input.Count,
	}, listquery.Options{})
	if err != nil {
		return nil, mapListQueryError(err)
	}

	services, err := h.useCase.List(ctx, catalogusecase.ServiceListFilter{
		Search: filter.Search,
		Offset: filter.Offset,
		Count:  filter.Count,
		Type:   serviceType,
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]BaseServiceResponse, 0, len(services))
	for _, item := range services {
		items = append(items, mapBaseResponse(item))
	}

	return &ListOutput{Body: items}, nil
}
