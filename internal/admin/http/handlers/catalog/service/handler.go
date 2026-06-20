package service

import catalogusecase "github.com/cheesy008/ffbc-backend/internal/catalog/use_case"

type Handler struct {
	useCase *catalogusecase.ServiceUseCase
}

func New(useCase *catalogusecase.ServiceUseCase) *Handler {
	return &Handler{useCase: useCase}
}
