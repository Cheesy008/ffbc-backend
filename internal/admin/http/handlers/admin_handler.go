package handlers

import (
	"context"

	admincontext "github.com/cheesy008/ffbc-backend/internal/admin/http/context"
)

type MeOutput struct {
	Body MeResponse
}
type MeResponse struct {
	ID    int64  `json:"id" example:"1"`
	Email string `json:"email" example:"admin@example.com"`
}

type AdminHandler struct {
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

type MeInput struct{}

func (h *AdminHandler) Me(ctx context.Context, _ *MeInput) (*MeOutput, error) {
	currentAdmin := admincontext.MustCurrentAdmin(ctx)

	return &MeOutput{
		Body: MeResponse{
			ID:    currentAdmin.ID,
			Email: currentAdmin.Email,
		},
	}, nil
}
