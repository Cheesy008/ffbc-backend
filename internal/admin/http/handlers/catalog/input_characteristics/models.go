package inputcharacteristics

import (
	"github.com/cheesy008/ffbc-backend/internal/shared/optional"
)

type InputCharacteristicCreateRequest struct {
	Name string `json:"name" required:"true" minLength:"1" maxLength:"250" example:"Длина изделия"`
	Type string `json:"type" required:"true" enum:"number,measuring,text,boolean,media" example:"number"`
}

type InputCharacteristicPatchRequest struct {
	Name        optional.Optional[string] `json:"name" required:"false" minLength:"1" maxLength:"250" example:"Длина изделия"`
	Type        optional.Optional[string] `json:"type" required:"false" enum:"number,measuring,text,boolean,media" example:"number"`
	TemplateIDs optional.Optional[*[]int] `json:"template_ids" required:"false"`
}

type InputCharacteristicBulkResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Длина изделия"`
	Type string `json:"type" example:"number"`
}

type InputCharacteristicListItemResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Длина изделия"`
	Type string `json:"type" example:"number"`
}

type InputCharacteristicResponse struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Длина изделия"`
	Type        string `json:"type" example:"number"`
	TemplateIDs []int  `json:"template_ids"`
	CreatedAt   string `json:"created_at" example:"2026-05-23T12:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2026-05-23T12:00:00Z"`
}
