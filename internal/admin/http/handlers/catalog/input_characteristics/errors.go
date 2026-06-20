package inputcharacteristics

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	catalogdomain "github.com/cheesy008/ffbc-backend/internal/catalog/domain"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, catalogdomain.ErrInputCharacteristicNotFound):
		return huma.Error404NotFound("input characteristic not found")
	case errors.Is(err, catalogdomain.ErrInputCharacteristicsTemplateNotFound):
		return huma.Error400BadRequest("input characteristics template not found")
	case errors.Is(err, catalogdomain.ErrEmptyPatch):
		return huma.Error400BadRequest("empty patch")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}

func mapListQueryError(err error) error {
	return huma.Error400BadRequest(err.Error())
}
