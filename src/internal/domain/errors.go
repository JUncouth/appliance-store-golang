package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource conflict")
	ErrValidation        = errors.New("validation failed")
	ErrInsufficientStock = errors.New("insufficient stock")
)
