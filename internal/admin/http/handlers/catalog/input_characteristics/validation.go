package inputcharacteristics

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
	catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"
)

func createInput(request InputCharacteristicCreateRequest) (catalogusecase.InputCharacteristicsInput, error) {
	name, err := normalizeName(request.Name)
	if err != nil {
		return catalogusecase.InputCharacteristicsInput{}, err
	}

	characteristicType, err := parseType(request.Type)
	if err != nil {
		return catalogusecase.InputCharacteristicsInput{}, err
	}

	return catalogusecase.InputCharacteristicsInput{
		Name: name,
		Type: characteristicType,
	}, nil
}

func patchInput(
	id int,
	request InputCharacteristicPatchRequest,
) (catalogusecase.InputCharacteristicsPatchInput, error) {
	input := catalogusecase.InputCharacteristicsPatchInput{ID: id}

	if request.Name.Set {
		name, err := normalizeName(request.Name.Value)
		if err != nil {
			return catalogusecase.InputCharacteristicsPatchInput{}, err
		}
		input.Name.Set = true
		input.Name.Value = name
	}

	if request.Type.Set {
		characteristicType, err := parseType(request.Type.Value)
		if err != nil {
			return catalogusecase.InputCharacteristicsPatchInput{}, err
		}
		input.Type.Set = true
		input.Type.Value = characteristicType
	}

	if request.TemplateIDs.Set {
		input.TemplateIDList.Set = true
		if request.TemplateIDs.Value != nil {
			templateIDs, err := normalizeTemplateIDs(*request.TemplateIDs.Value)
			if err != nil {
				return catalogusecase.InputCharacteristicsPatchInput{}, err
			}
			input.TemplateIDList.Value = &templateIDs
		}
	}

	return input, nil
}

func normalizeName(value string) (string, error) {
	name := strings.TrimSpace(value)
	if name == "" || len([]rune(name)) > 250 {
		return "", huma.Error400BadRequest("invalid input characteristic name")
	}
	return name, nil
}

func parseType(value string) (catalogdomain.InputCharacteristicType, error) {
	characteristicType := catalogdomain.InputCharacteristicType(value)
	switch characteristicType {
	case catalogdomain.InputCharacteristicTypeNumber,
		catalogdomain.InputCharacteristicTypeMeasuring,
		catalogdomain.InputCharacteristicTypeText,
		catalogdomain.InputCharacteristicTypeBoolean,
		catalogdomain.InputCharacteristicTypeMedia:
		return characteristicType, nil
	default:
		return "", huma.Error400BadRequest("invalid input characteristic type")
	}
}

func normalizeTemplateIDs(values []int) ([]int, error) {
	result := make([]int, 0, len(values))
	seen := make(map[int]struct{}, len(values))
	for _, value := range values {
		if value <= 0 {
			return nil, huma.Error400BadRequest("invalid template id")
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result, nil
}
