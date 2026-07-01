package models

import "errors"

var(
	ErrNotFound    = errors.New("resource not found")
	ErrDuplicate   = errors.New("resource already exists")
	ErrValidation  = errors.New("validation failed")
	ErrOrgNotFound = errors.New("organization does not exist")
)