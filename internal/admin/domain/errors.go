package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminInactive      = errors.New("admin inactive")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminAlreadyExists = errors.New("admin email already exists")

	ErrorSessionCreation = errors.New("error creating session")
	ErrorSessionNotFound = errors.New("error session not found")
	ErrorSessionRevoke   = errors.New("error session revoke")
	ErrorSessionExpired  = errors.New("error session expired")
)
