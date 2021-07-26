package router

import "gopkg.in/go-playground/validator.v9"

// NewValidator -  returns new Validator
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validator - Place holder for Validator
type Validator struct {
	validator *validator.Validate
}

// Validate
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
