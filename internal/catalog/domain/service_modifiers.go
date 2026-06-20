package domain

import "time"

type ServiceModifier struct {
	ID            int
	ServiceID     int
	Name          string
	SelectionType ServiceModifierSelectionType
	SortOrder     *int
	IsRequired    bool
	Values        []ServiceModifierValue
}

type ServiceModifierValue struct {
	ID                int
	Name              string
	ServiceModifierID int
	AdditionalPrice   Decimal
	IsActive          bool
	SortOrder         *int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
