package domain

import "time"

type AdminUser struct {
	ID           int64
	Email        string
	PasswordHash string
	DisplayName  *string
	IsActive     bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u AdminUser) CanLogin() bool {
	return u.IsActive
}
