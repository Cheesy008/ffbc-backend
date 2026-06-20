package inputcharacteristics

import catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"

type Handler struct {
	useCase *catalogusecase.InputCharacteristicsUseCase
}

func New(useCase *catalogusecase.InputCharacteristicsUseCase) *Handler {
	return &Handler{useCase: useCase}
}
