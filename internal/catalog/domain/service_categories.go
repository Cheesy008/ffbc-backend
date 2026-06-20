package domain

import "time"

type ServiceCategory struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ServiceCategoryLink struct {
	ServiceID         int
	ServiceCategoryID int
}
