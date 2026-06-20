package domain

import "time"

type InputCharacteristicTemplate struct {
	ID                   int
	Name                 string
	Description          *string
	InputCharacteristics []InputCharacteristic
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type InputCharacteristic struct {
	ID          int
	Name        string
	Type        InputCharacteristicType
	TemplateIDs []int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InputCharacteristicTemplateItem struct {
	TemplateID            int
	InputCharacteristicID int
}

type ServiceInputCharacteristic struct {
	InputCharacteristic InputCharacteristic
	ServiceID           int
	IsRequired          bool
	SortOrder           *int
}
