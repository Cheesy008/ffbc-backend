package domain

import "time"

type Service struct {
	ID                   int
	Name                 string
	BasePrice            Decimal
	Description          *string
	Type                 ServiceType
	Status               ServiceStatus
	InputCharacteristics []ServiceInputCharacteristic
	Modifiers            []ServiceModifier
	Categories           []ServiceCategory
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ServiceListItem struct {
	ID          int
	Name        string
	BasePrice   Decimal
	Description *string
	Type        ServiceType
}

type ServiceMediaFile struct {
	ID        int
	ServiceID int
	FileURL   string
	SortOrder *int
	CreatedAt time.Time
	UpdatedAt time.Time
}
