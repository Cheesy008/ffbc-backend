package service

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shopspring/decimal"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
)

func sellingCreateInput(request SellingServiceCreateRequest) (catalogusecase.ServiceInput, error) {
	name, err := normalizeName(request.Name)
	if err != nil {
		return catalogusecase.ServiceInput{}, err
	}
	if err := validateDescription(request.Description); err != nil {
		return catalogusecase.ServiceInput{}, err
	}

	return catalogusecase.ServiceInput{
		Name:        name,
		BasePrice:   decimal.NewFromFloat(request.BasePrice),
		Description: request.Description,
		Type:        catalogdomain.ServiceTypeSelling,
		Status:      catalogdomain.ServiceStatusBaseInfoRequired,
	}, nil
}

func creationCreateInput(request CreationServiceCreateRequest) (catalogusecase.ServiceInput, error) {
	input, err := sellingCreateInput(SellingServiceCreateRequest{
		Name:        request.Name,
		BasePrice:   request.BasePrice,
		Description: request.Description,
	})
	if err != nil {
		return catalogusecase.ServiceInput{}, err
	}
	input.Type = catalogdomain.ServiceTypeCreation

	characteristics := make([]catalogusecase.ServiceInputCharacteristicInput, 0, len(request.InputCharacteristics))
	seen := make(map[int]struct{}, len(request.InputCharacteristics))
	for _, item := range request.InputCharacteristics {
		id := item.InputCharacteristic
		if item.SortOrder != nil && *item.SortOrder < 0 {
			return catalogusecase.ServiceInput{}, huma.Error400BadRequest("invalid sort order")
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		characteristics = append(characteristics, catalogusecase.ServiceInputCharacteristicInput{
			InputCharacteristicID: id,
			IsRequired:            item.IsRequired,
			SortOrder:             item.SortOrder,
		})
	}
	input.InputCharacteristics = characteristics

	return input, nil
}

func sellingPatchInput(
	id int,
	request SellingServicePatchRequest,
) (catalogusecase.ServicePatchInput, error) {
	input := catalogusecase.ServicePatchInput{
		ID:   id,
		Type: catalogdomain.ServiceTypeSelling,
	}

	if request.Name.Set {
		name, err := normalizeName(request.Name.Value)
		if err != nil {
			return catalogusecase.ServicePatchInput{}, err
		}
		input.Name.Set = true
		input.Name.Value = name
	}
	if request.BasePrice.Set {
		if request.BasePrice.Value < 0 {
			return catalogusecase.ServicePatchInput{}, huma.Error400BadRequest("invalid base price")
		}
		input.BasePrice.Set = true
		input.BasePrice.Value = decimal.NewFromFloat(request.BasePrice.Value)
	}
	if request.Description.Set {
		if err := validateDescription(request.Description.Value); err != nil {
			return catalogusecase.ServicePatchInput{}, err
		}
		input.Description = request.Description
	}

	return input, nil
}

func creationPatchInput(
	id int,
	request CreationServicePatchRequest,
) (catalogusecase.ServicePatchInput, error) {
	input, err := sellingPatchInput(id, SellingServicePatchRequest{
		Name:        request.Name,
		BasePrice:   request.BasePrice,
		Description: request.Description,
	})
	if err != nil {
		return catalogusecase.ServicePatchInput{}, err
	}
	input.Type = catalogdomain.ServiceTypeCreation

	if request.InputCharacteristics.Set {
		if request.InputCharacteristics.Value == nil {
			return catalogusecase.ServicePatchInput{},
				huma.Error400BadRequest("input characteristics must be an array")
		}
		ids, err := parseIDs(request.InputCharacteristics.Value)
		if err != nil {
			return catalogusecase.ServicePatchInput{}, err
		}
		input.InputCharacteristicIDs.Set = true
		input.InputCharacteristicIDs.Value = ids
	}

	return input, nil
}

func normalizeName(value string) (string, error) {
	name := strings.TrimSpace(value)
	if name == "" || len([]rune(name)) > 250 {
		return "", huma.Error400BadRequest("invalid service name")
	}
	return name, nil
}

func validateDescription(value *string) error {
	if value != nil && len([]rune(*value)) > 1500 {
		return huma.Error400BadRequest("invalid description")
	}
	return nil
}

func parseIDs(values []int) ([]int, error) {
	result := make([]int, 0, len(values))
	seen := make(map[int]struct{}, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result, nil
}
