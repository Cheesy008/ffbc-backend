package catalog

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
	"github.com/cheesy008/ffbc-backend/internal/shared/listquery"
)

type ServiceCategoryHandler struct {
	categoryUseCase *catalogusecase.ServiceCategoryUseCase
}

func NewServiceCategoryHandler(categoryUseCase *catalogusecase.ServiceCategoryUseCase) *ServiceCategoryHandler {
	return &ServiceCategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

type ServiceCategoryCreateInput struct {
	Body ServiceCategoryRequest
}

type ServiceCategoryUpdateInput struct {
	ID   int `path:"id" minimum:"1" example:"1"`
	Body ServiceCategoryRequest
}

type ServiceCategoryListInput struct {
	Search    string `query:"search" required:"false" example:"Одежда"`
	Offset    int    `query:"offset" minimum:"0" default:"0" example:"0"`
	Count     int    `query:"count" minimum:"1" maximum:"100" default:"20" example:"20"`
	SortBy    string `query:"sortBy" enum:"name" default:"name" example:"name"`
	SortOrder string `query:"sortOrder" enum:"asc,desc" default:"asc" example:"asc"`
}

type ServiceCategoryDeleteInput struct {
	ID int `path:"id" minimum:"1" example:"1"`
}

type ServiceCategoryRequest struct {
	Name string `json:"name" required:"true" minLength:"1" maxLength:"250"`
}

type ServiceCategoryOutput struct {
	Body ServiceCategoryResponse
}

type ServiceCategoryListOutput struct {
	Body []ServiceCategoryResponse
}

type ServiceCategoryResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Штаны"`
}

type ServiceCategoryDeleteOutput struct{}

func (h *ServiceCategoryHandler) Create(
	ctx context.Context,
	input *ServiceCategoryCreateInput,
) (*ServiceCategoryOutput, error) {
	category, err := h.categoryUseCase.Create(ctx, catalogusecase.ServiceCategoryInput{
		Name: input.Body.Name,
	})
	if err != nil {
		return nil, mapServiceCategoryError(err)
	}

	return &ServiceCategoryOutput{Body: mapServiceCategoryResponse(category)}, nil
}

func (h *ServiceCategoryHandler) Update(
	ctx context.Context,
	input *ServiceCategoryUpdateInput,
) (*ServiceCategoryOutput, error) {
	category, err := h.categoryUseCase.Update(ctx, catalogusecase.ServiceCategoryInput{
		ID:   input.ID,
		Name: input.Body.Name,
	})
	if err != nil {
		return nil, mapServiceCategoryError(err)
	}

	return &ServiceCategoryOutput{Body: mapServiceCategoryResponse(category)}, nil
}

func (h *ServiceCategoryHandler) List(
	ctx context.Context,
	input *ServiceCategoryListInput,
) (*ServiceCategoryListOutput, error) {
	filter, err := listquery.NormalizeNameSorted(listquery.Input{
		Search:    input.Search,
		Offset:    input.Offset,
		Count:     input.Count,
		SortBy:    input.SortBy,
		SortOrder: input.SortOrder,
	})
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	categories, err := h.categoryUseCase.List(ctx, catalogusecase.ServiceCategoryListFilter{
		Search:    filter.Search,
		Offset:    filter.Offset,
		Count:     filter.Count,
		SortBy:    filter.SortBy,
		SortOrder: filter.SortOrder,
	})
	if err != nil {
		return nil, mapServiceCategoryError(err)
	}

	items := make([]ServiceCategoryResponse, 0, len(categories))
	for _, category := range categories {
		items = append(items, mapServiceCategoryResponse(category))
	}

	return &ServiceCategoryListOutput{Body: items}, nil
}

func (h *ServiceCategoryHandler) Delete(
	ctx context.Context,
	input *ServiceCategoryDeleteInput,
) (*ServiceCategoryDeleteOutput, error) {
	if err := h.categoryUseCase.Delete(ctx, input.ID); err != nil {
		if errors.Is(err, catalogdomain.ErrServiceCategoryNotFound) {
			return nil, huma.Error400BadRequest("service category not found")
		}
		return nil, mapServiceCategoryError(err)
	}

	return &ServiceCategoryDeleteOutput{}, nil
}

func mapServiceCategoryError(err error) error {
	switch {
	case errors.Is(err, catalogdomain.ErrServiceCategoryNotFound):
		return huma.Error404NotFound("service category not found")
	case errors.Is(err, catalogdomain.ErrServiceCategoryAlreadyExists):
		return huma.Error409Conflict("service category already exists")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}

func mapServiceCategoryResponse(category catalogdomain.ServiceCategory) ServiceCategoryResponse {
	return ServiceCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}
