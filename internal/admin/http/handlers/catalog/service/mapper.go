package service

import (
	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

func mapBaseResponse(item catalogdomain.ServiceListItem) BaseServiceResponse {
	return BaseServiceResponse{
		ID:          item.ID,
		Name:        item.Name,
		BasePrice:   item.BasePrice.InexactFloat64(),
		Description: item.Description,
		Media:       []ServiceMediaResponse{},
	}
}

func mapSellingResponse(item catalogdomain.Service) SellingServiceResponse {
	return SellingServiceResponse{
		ID:          item.ID,
		Name:        item.Name,
		BasePrice:   item.BasePrice.InexactFloat64(),
		Description: item.Description,
		Status:      string(item.Status),
		Modifiers:   mapModifiers(item.Modifiers),
		Categories:  mapCategories(item.Categories),
		Media:       []ServiceMediaResponse{},
	}
}

func mapCreationResponse(item catalogdomain.Service) CreationServiceResponse {
	return CreationServiceResponse{
		ID:                   item.ID,
		Name:                 item.Name,
		BasePrice:            item.BasePrice.InexactFloat64(),
		Description:          item.Description,
		Status:               string(item.Status),
		InputCharacteristics: mapInputCharacteristics(item.InputCharacteristics),
		Modifiers:            mapModifiers(item.Modifiers),
		Categories:           mapCategories(item.Categories),
		Media:                []ServiceMediaResponse{},
	}
}

func mapInputCharacteristics(
	items []catalogdomain.ServiceInputCharacteristic,
) []ServiceInputCharacteristicResponse {
	result := make([]ServiceInputCharacteristicResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ServiceInputCharacteristicResponse{
			InputCharacteristic: ServiceInputCharacteristicItemResponse{
				ID:   item.InputCharacteristic.ID,
				Name: item.InputCharacteristic.Name,
				Type: string(item.InputCharacteristic.Type),
			},
			IsRequired: item.IsRequired,
			SortOrder:  item.SortOrder,
		})
	}
	return result
}

func mapCategories(items []catalogdomain.ServiceCategory) []ServiceCategoryItemResponse {
	result := make([]ServiceCategoryItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ServiceCategoryItemResponse{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	return result
}

func mapModifiers(items []catalogdomain.ServiceModifier) []ServiceModifierResponse {
	result := make([]ServiceModifierResponse, 0, len(items))
	for _, item := range items {
		values := make([]ServiceModifierValueResponse, 0, len(item.Values))
		for _, value := range item.Values {
			values = append(values, ServiceModifierValueResponse{
				ID:              value.ID,
				Name:            value.Name,
				AdditionalPrice: value.AdditionalPrice.InexactFloat64(),
				IsActive:        value.IsActive,
				SortOrder:       value.SortOrder,
			})
		}
		result = append(result, ServiceModifierResponse{
			ID:            item.ID,
			ServiceID:     item.ServiceID,
			Name:          item.Name,
			SelectionType: string(item.SelectionType),
			SortOrder:     item.SortOrder,
			IsRequired:    item.IsRequired,
			Values:        values,
		})
	}
	return result
}
