// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

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

// Validate - Validates given interface
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
