package models

import (
	"fmt"

	"errors"
)

type ValidationError struct {
	Fields []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("required fields are empty: %v", e.Fields)
}

func NewValidationError(fields ...string) error {
	return &ValidationError{Fields: fields}
}

var (
	ErrStoreIDRequired = errors.New("store ID is required")
)
