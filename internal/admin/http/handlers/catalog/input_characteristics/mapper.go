package inputcharacteristics

import catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"

func mapBulkResponse(characteristic catalogdomain.InputCharacteristic) InputCharacteristicBulkResponse {
	return InputCharacteristicBulkResponse{
		ID:   characteristic.ID,
		Name: characteristic.Name,
		Type: string(characteristic.Type),
	}
}

func mapListItemResponse(characteristic catalogdomain.InputCharacteristic) InputCharacteristicListItemResponse {
	return InputCharacteristicListItemResponse{
		ID:   characteristic.ID,
		Name: characteristic.Name,
		Type: string(characteristic.Type),
	}
}

func mapResponse(characteristic catalogdomain.InputCharacteristic) InputCharacteristicResponse {
	return InputCharacteristicResponse{
		ID:          characteristic.ID,
		Name:        characteristic.Name,
		Type:        string(characteristic.Type),
		TemplateIDs: characteristic.TemplateIDs,
		CreatedAt:   characteristic.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   characteristic.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
