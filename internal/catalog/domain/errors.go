package domain

import "errors"

var (
	ErrServiceCategoryNotFound              = errors.New("service category not found")
	ErrServiceCategoryAlreadyExists         = errors.New("service category already exists")
	ErrServiceNotFound                      = errors.New("service not found")
	ErrServiceTypeMismatch                  = errors.New("service type mismatch")
	ErrInputCharacteristicsTemplateNotFound = errors.New("input characteristics template not found")
	ErrInputCharacteristicNotFound          = errors.New("input characteristic not found")
	ErrEmptyPatch                           = errors.New("empty patch")
)
