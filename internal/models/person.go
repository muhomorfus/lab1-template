package models

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidData = errors.New("invalid data")
	ErrNotFound    = errors.New("not found")
)

type Person struct {
	ID      int
	Name    string `validate:"required"`
	Address *string
	Age     *int
	Work    *string
}

func (p *Person) Validate() error {
	if err := validator.New().Struct(p); err != nil {
		return fmt.Errorf("validate person: %w (%w)", err, ErrInvalidData)
	}

	return nil
}
