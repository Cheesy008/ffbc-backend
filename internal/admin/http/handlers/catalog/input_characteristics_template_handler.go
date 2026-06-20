package catalog

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cheesy008/ffbc-backend/internal/shared/optional"
	"github.com/danielgtaylor/huma/v2"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
	"github.com/cheesy008/ffbc-backend/internal/shared/listquery"
)

type InputCharacteristicsTemplateHandler struct {
	templateUseCase *catalogusecase.InputCharacteristicsTemplateUseCase
}

func NewInputCharacteristicsTemplateHandler(
	templateUseCase *catalogusecase.InputCharacteristicsTemplateUseCase,
) *InputCharacteristicsTemplateHandler {
	return &InputCharacteristicsTemplateHandler{
		templateUseCase: templateUseCase,
	}
}

type InputCharacteristicsTemplateCreateInput struct {
	Body InputCharacteristicsTemplateCreateRequest
}

type InputCharacteristicsTemplatePatchInput struct {
	ID   int `path:"id" minimum:"1" example:"1"`
	Body InputCharacteristicsTemplatePatchRequest
}

type InputCharacteristicsTemplateDetailInput struct {
	ID int `path:"id" minimum:"1" example:"1"`
}

type InputCharacteristicsTemplateDeleteInput struct {
	ID int `path:"id" minimum:"1" example:"1"`
}

type InputCharacteristicsTemplateDeleteOutput struct{}

type InputCharacteristicsTemplateListInput struct {
	Search               string   `query:"search" required:"false" example:"Мерки"`
	InputCharacteristics []string `query:"inputCharacteristics" required:"false"`
	Offset               int      `query:"offset" minimum:"0" default:"0" example:"0"`
	Count                int      `query:"count" minimum:"1" maximum:"100" default:"20" example:"20"`
	SortBy               string   `query:"sortBy" enum:"name" default:"name" example:"name"`
	SortOrder            string   `query:"sortOrder" enum:"asc,desc" default:"asc" example:"asc"`
}

type InputCharacteristicsTemplateCreateRequest struct {
	Name                 string   `json:"name" required:"true" minLength:"1" maxLength:"250" example:"Мерки для одежды"`
	InputCharacteristics []string `json:"inputCharacteristics" required:"true"`
	Description          *string  `json:"description" required:"true" maxLength:"1500"`
}

type InputCharacteristicsTemplatePatchRequest struct {
	Name                 optional.Optional[string]   `json:"name" required:"false" minLength:"1" maxLength:"250" example:"Мерки для одежды"`
	InputCharacteristics optional.Optional[[]string] `json:"inputCharacteristics" required:"false"`
	Description          optional.Optional[*string]  `json:"description" required:"false"`
}

type InputCharacteristicsTemplateDetailOutput struct {
	Body InputCharacteristicsTemplateDetailResponse
}

type InputCharacteristicsTemplateListOutput struct {
	Body []InputCharacteristicsTemplateShortResponse
}

type InputCharacteristicsTemplateShortResponse struct {
	ID          int     `json:"id" example:"1"`
	Name        string  `json:"name" example:"Мерки для одежды"`
	Description *string `json:"description"`
}

type InputCharacteristicsTemplateDetailResponse struct {
	ID                   string                                    `json:"id" example:"1"`
	Name                 string                                    `json:"name" example:"Мерки для одежды"`
	InputCharacteristics []InputCharacteristicTemplateItemResponse `json:"inputCharacteristics"`
	Description          *string                                   `json:"description"`
}

type InputCharacteristicTemplateItemResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Длина изделия"`
	Type string `json:"type" example:"number"`
}

func (h *InputCharacteristicsTemplateHandler) Create(
	ctx context.Context,
	input *InputCharacteristicsTemplateCreateInput,
) (*InputCharacteristicsTemplateDetailOutput, error) {
	templateInput, err := inputCharacteristicsTemplateInput(input.Body)
	if err != nil {
		return nil, err
	}

	template, err := h.templateUseCase.Create(ctx, templateInput)
	if err != nil {
		return nil, mapInputCharacteristicsTemplateError(err)
	}

	return &InputCharacteristicsTemplateDetailOutput{
		Body: mapInputCharacteristicsTemplateDetailResponse(template),
	}, nil
}

func (h *InputCharacteristicsTemplateHandler) Patch(
	ctx context.Context,
	input *InputCharacteristicsTemplatePatchInput,
) (*InputCharacteristicsTemplateDetailOutput, error) {
	templateInput, err := inputCharacteristicsTemplatePatchInput(input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	template, err := h.templateUseCase.Patch(ctx, templateInput)
	if err != nil {
		return nil, mapInputCharacteristicsTemplateError(err)
	}

	return &InputCharacteristicsTemplateDetailOutput{
		Body: mapInputCharacteristicsTemplateDetailResponse(template),
	}, nil
}

func (h *InputCharacteristicsTemplateHandler) Detail(
	ctx context.Context,
	input *InputCharacteristicsTemplateDetailInput,
) (*InputCharacteristicsTemplateDetailOutput, error) {
	template, err := h.templateUseCase.GetByID(ctx, input.ID)
	if err != nil {
		return nil, mapInputCharacteristicsTemplateError(err)
	}

	return &InputCharacteristicsTemplateDetailOutput{
		Body: mapInputCharacteristicsTemplateDetailResponse(template),
	}, nil
}

func (h *InputCharacteristicsTemplateHandler) Delete(
	ctx context.Context,
	input *InputCharacteristicsTemplateDeleteInput,
) (*InputCharacteristicsTemplateDeleteOutput, error) {
	if err := h.templateUseCase.Delete(ctx, input.ID); err != nil {
		if errors.Is(err, catalogdomain.ErrInputCharacteristicsTemplateNotFound) {
			return nil, huma.Error400BadRequest("input characteristics template not found")
		}
		return nil, mapInputCharacteristicsTemplateError(err)
	}

	return &InputCharacteristicsTemplateDeleteOutput{}, nil
}

func (h *InputCharacteristicsTemplateHandler) List(
	ctx context.Context,
	input *InputCharacteristicsTemplateListInput,
) (*InputCharacteristicsTemplateListOutput, error) {
	inputCharacteristicIDs, err := parseInputCharacteristicIDs(input.InputCharacteristics)
	if err != nil {
		return nil, err
	}

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

	templates, err := h.templateUseCase.List(ctx, catalogusecase.InputCharacteristicsTemplateListFilter{
		Search:                 filter.Search,
		InputCharacteristicIDs: inputCharacteristicIDs,
		Offset:                 filter.Offset,
		Count:                  filter.Count,
		SortBy:                 filter.SortBy,
		SortOrder:              filter.SortOrder,
	})
	if err != nil {
		return nil, mapInputCharacteristicsTemplateError(err)
	}

	items := make([]InputCharacteristicsTemplateShortResponse, 0, len(templates))
	for _, template := range templates {
		items = append(items, mapInputCharacteristicsTemplateShortResponse(template))
	}

	return &InputCharacteristicsTemplateListOutput{Body: items}, nil
}

func inputCharacteristicsTemplateInput(
	request InputCharacteristicsTemplateCreateRequest,
) (catalogusecase.InputCharacteristicsTemplateInput, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" || len([]rune(name)) > 250 {
		return catalogusecase.InputCharacteristicsTemplateInput{}, huma.Error400BadRequest("invalid input characteristics template name")
	}
	if request.Description != nil && len([]rune(*request.Description)) > 1500 {
		return catalogusecase.InputCharacteristicsTemplateInput{}, huma.Error400BadRequest("invalid description")
	}

	inputCharacteristicIDs, err := parseInputCharacteristicIDs(request.InputCharacteristics)
	if err != nil {
		return catalogusecase.InputCharacteristicsTemplateInput{}, err
	}

	return catalogusecase.InputCharacteristicsTemplateInput{
		Name:                   name,
		Description:            request.Description,
		InputCharacteristicIDs: inputCharacteristicIDs,
	}, nil
}

func parseInputCharacteristicIDs(values []string) ([]int, error) {
	result := make([]int, 0, len(values))
	seen := make(map[int]struct{}, len(values))

	for _, value := range values {
		id, err := strconv.Atoi(value)
		if err != nil || id <= 0 {
			return nil, huma.Error400BadRequest("invalid input characteristic id")
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result, nil
}

func inputCharacteristicsTemplatePatchInput(
	id int,
	request InputCharacteristicsTemplatePatchRequest,
) (catalogusecase.InputCharacteristicsTemplatePatchInput, error) {
	input := catalogusecase.InputCharacteristicsTemplatePatchInput{ID: id}

	if request.Name.Set {
		name := strings.TrimSpace(request.Name.Value)
		if name == "" || len([]rune(name)) > 250 {
			return catalogusecase.InputCharacteristicsTemplatePatchInput{}, huma.Error400BadRequest("invalid input characteristics template name")
		}
		input.Name.Set = true
		input.Name.Value = name
	}

	if request.Description.Set {
		if request.Description.Value != nil && len([]rune(*request.Description.Value)) > 1500 {
			return catalogusecase.InputCharacteristicsTemplatePatchInput{}, huma.Error400BadRequest("invalid description")
		}
		input.Description = request.Description
	}

	if request.InputCharacteristics.Set {
		if request.InputCharacteristics.Value == nil {
			return catalogusecase.InputCharacteristicsTemplatePatchInput{},
				huma.Error400BadRequest("input characteristics must be an array")
		}

		inputCharacteristicIDs, err := parseInputCharacteristicIDs(request.InputCharacteristics.Value)
		if err != nil {
			return catalogusecase.InputCharacteristicsTemplatePatchInput{}, err
		}
		input.InputCharacteristicIDs.Set = true
		input.InputCharacteristicIDs.Value = inputCharacteristicIDs
	}

	return input, nil
}

func mapInputCharacteristicsTemplateError(err error) error {
	switch {
	case errors.Is(err, catalogdomain.ErrInputCharacteristicsTemplateNotFound):
		return huma.Error404NotFound("input characteristics template not found")
	case errors.Is(err, catalogdomain.ErrInputCharacteristicNotFound):
		return huma.Error400BadRequest("input characteristic not found")
	case errors.Is(err, catalogdomain.ErrEmptyPatch):
		return huma.Error400BadRequest("empty patch")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}

func mapInputCharacteristicsTemplateShortResponse(
	template catalogdomain.InputCharacteristicTemplate,
) InputCharacteristicsTemplateShortResponse {
	return InputCharacteristicsTemplateShortResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
	}
}

func mapInputCharacteristicsTemplateDetailResponse(
	template catalogdomain.InputCharacteristicTemplate,
) InputCharacteristicsTemplateDetailResponse {
	inputCharacteristics := make([]InputCharacteristicTemplateItemResponse, 0, len(template.InputCharacteristics))
	for _, characteristic := range template.InputCharacteristics {
		inputCharacteristics = append(inputCharacteristics, InputCharacteristicTemplateItemResponse{
			ID:   characteristic.ID,
			Name: characteristic.Name,
			Type: string(characteristic.Type),
		})
	}

	return InputCharacteristicsTemplateDetailResponse{
		ID:                   strconv.Itoa(template.ID),
		Name:                 template.Name,
		InputCharacteristics: inputCharacteristics,
		Description:          template.Description,
	}
}
