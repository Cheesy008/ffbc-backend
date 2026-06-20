package domain

import "time"

type ServiceProduct struct {
	ID                  int
	ServiceID           int
	Name                string
	AutoCalculatedPrice Decimal
	ConfirmedPrice      Decimal
	Comment             *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ServiceProductInputCharacteristicValue struct {
	ID                    int
	InputCharacteristicID int
	ServiceProductID      int
	NumberValue           *Decimal
	MeasuringValue        *Decimal
	TextValue             *string
	BooleanValue          *bool
	MediaFile             *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type ServiceProductModifierValue struct {
	ServiceProductID int
	ModifierValueID  int
}
