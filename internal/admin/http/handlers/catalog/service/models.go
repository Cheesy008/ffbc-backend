package service

import (
	"github.com/cheesy008/ffbc-backend/internal/shared/optional"
)

type ListInput struct {
	Search string `query:"search" required:"false" example:"Платье"`
	Offset int    `query:"offset" minimum:"0" default:"0" example:"0"`
	Count  int    `query:"count" minimum:"1" maximum:"100" default:"20" example:"20"`
}

type DetailInput struct {
	ID int `path:"id" minimum:"1" example:"1"`
}

type DeleteInput struct {
	ID int `path:"id" minimum:"1" example:"1"`
}

type SellingCreateInput struct {
	Body SellingServiceCreateRequest
}

type CreationCreateInput struct {
	Body CreationServiceCreateRequest
}

type SellingPatchInput struct {
	ID   int `path:"id" minimum:"1" example:"1"`
	Body SellingServicePatchRequest
}

type CreationPatchInput struct {
	ID   int `path:"id" minimum:"1" example:"1"`
	Body CreationServicePatchRequest
}

type SellingServiceCreateRequest struct {
	Name        string  `json:"name" required:"true" minLength:"1" maxLength:"250"`
	BasePrice   float64 `json:"basePrice" required:"true" minimum:"0"`
	Description *string `json:"description" required:"true"`
}

type CreationServiceCreateRequest struct {
	Name                 string                                    `json:"name" required:"true" minLength:"1" maxLength:"250"`
	BasePrice            float64                                   `json:"basePrice" required:"true" minimum:"0"`
	Description          *string                                   `json:"description" required:"true"`
	InputCharacteristics []ServiceInputCharacteristicCreateRequest `json:"inputCharacteristics" required:"true"`
}

type ServiceInputCharacteristicCreateRequest struct {
	InputCharacteristic int  `json:"inputCharacteristic" required:"true"`
	IsRequired          bool `json:"isRequired" required:"true"`
	SortOrder           *int `json:"sortOrder" required:"true" minimum:"0"`
}

type SellingServicePatchRequest struct {
	Name        optional.Optional[string]  `json:"name" required:"false" minLength:"1" maxLength:"250"`
	BasePrice   optional.Optional[float64] `json:"basePrice" required:"false" minimum:"0"`
	Description optional.Optional[*string] `json:"description" required:"false"`
}

type CreationServicePatchRequest struct {
	Name                 optional.Optional[string]  `json:"name" required:"false" minLength:"1" maxLength:"250"`
	BasePrice            optional.Optional[float64] `json:"basePrice" required:"false" minimum:"0"`
	Description          optional.Optional[*string] `json:"description" required:"false"`
	InputCharacteristics optional.Optional[[]int]   `json:"inputCharacteristics" required:"false"`
}

type ListOutput struct {
	Body []BaseServiceResponse
}

type SellingOutput struct {
	Body SellingServiceResponse
}

type CreationOutput struct {
	Body CreationServiceResponse
}

type DeleteOutput struct{}

type BaseServiceResponse struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	BasePrice   float64                `json:"basePrice"`
	Description *string                `json:"description"`
	Media       []ServiceMediaResponse `json:"media"`
}

type SellingServiceResponse struct {
	ID          int                           `json:"id"`
	Name        string                        `json:"name"`
	BasePrice   float64                       `json:"basePrice"`
	Description *string                       `json:"description"`
	Status      string                        `json:"status"`
	Modifiers   []ServiceModifierResponse     `json:"modifiers"`
	Categories  []ServiceCategoryItemResponse `json:"categories"`
	Media       []ServiceMediaResponse        `json:"media"`
}

type CreationServiceResponse struct {
	ID                   int                                  `json:"id"`
	Name                 string                               `json:"name"`
	BasePrice            float64                              `json:"basePrice"`
	Description          *string                              `json:"description"`
	Status               string                               `json:"status"`
	InputCharacteristics []ServiceInputCharacteristicResponse `json:"inputCharacteristics"`
	Modifiers            []ServiceModifierResponse            `json:"modifiers"`
	Categories           []ServiceCategoryItemResponse        `json:"categories"`
	Media                []ServiceMediaResponse               `json:"media"`
}

type ServiceInputCharacteristicResponse struct {
	InputCharacteristic ServiceInputCharacteristicItemResponse `json:"inputCharacteristic"`
	IsRequired          bool                                   `json:"isRequired"`
	SortOrder           *int                                   `json:"sortOrder"`
}

type ServiceInputCharacteristicItemResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ServiceCategoryItemResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ServiceModifierResponse struct {
	ID            int                            `json:"id"`
	ServiceID     int                            `json:"serviceId"`
	Name          string                         `json:"name"`
	SelectionType string                         `json:"selectionType"`
	SortOrder     *int                           `json:"sortOrder"`
	IsRequired    bool                           `json:"isRequired"`
	Values        []ServiceModifierValueResponse `json:"values"`
}

type ServiceModifierValueResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	AdditionalPrice float64 `json:"additionalPrice"`
	IsActive        bool    `json:"isActive"`
	SortOrder       *int    `json:"sortOrder"`
}

type ServiceMediaResponse struct {
	File      ServiceFileResponse `json:"file"`
	AltText   *string             `json:"altText"`
	Caption   *string             `json:"caption"`
	SortOrder *int                `json:"sortOrder"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
}

type ServiceFileResponse struct {
	ID         string `json:"id"`
	StorageKey string `json:"storageKey"`
	MimeType   string `json:"mimeType"`
	Size       int64  `json:"size" minimum:"0"`
	UpdatedAt  string `json:"updatedAt"`
}
