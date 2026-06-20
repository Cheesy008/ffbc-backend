package domain

import "github.com/shopspring/decimal"

type Decimal = decimal.Decimal

type InputCharacteristicType string

const (
	InputCharacteristicTypeNumber    InputCharacteristicType = "number"
	InputCharacteristicTypeMeasuring                         = "measuring"
	InputCharacteristicTypeText                              = "text"
	InputCharacteristicTypeBoolean                           = "boolean"
	InputCharacteristicTypeMedia                             = "media"
	InputCharacteristicTypeVideo                             = "video"
)

type ServiceModifierSelectionType string

const (
	ServiceModifierSelectionTypeSingle   ServiceModifierSelectionType = "single"
	ServiceModifierSelectionTypeMultiple                              = "multiple"
)

type ServiceType string

const (
	ServiceTypeCreation ServiceType = "creation"
	ServiceTypeSelling              = "selling"
)

type ServiceStatus string

const (
	ServiceStatusBaseInfoRequired ServiceStatus = "base_info_required"
	ServiceStatusMediaRequired                  = "media_required"
	ServiceStatusActive                         = "active"
	ServiceStatusDisabled                       = "disabled"
)
