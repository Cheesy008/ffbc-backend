package domain

import "time"

type AdminSession struct {
	ID          int64
	AdminUserID int64
	TokenHash   string
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
	LastUsedAt  *time.Time
}

func (s AdminSession) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && s.ExpiresAt.After(now)
}
