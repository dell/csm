// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package store

import (
	"errors"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

// UserStoreInterface is used to define the interface for persisting User
//go:generate mockgen -destination=mocks/user_store_interface.go -package=mocks github.com/dell/csm-deployment/store UserStoreInterface
type UserStoreInterface interface {
	GetByID(uint) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	Create(*model.User) error
	Update(*model.User) error
}

// UserStore - Placeholder for User store
type UserStore struct {
	db *gorm.DB
}

// NewUserStore - returns an instance of UserStore in db
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// GetByID - Method to get user by Id
func (us *UserStore) GetByID(id uint) (*model.User, error) {
	var m model.User
	if err := us.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// GetByUsername - Method to get User by name
func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := us.db.Where(&model.User{Username: username}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// Create - Method to create User
func (us *UserStore) Create(u *model.User) (err error) {
	return us.db.Create(u).Error
}

// Update - Method to update User
func (us *UserStore) Update(u *model.User) error {
	return us.db.Model(u).Updates(u).Error
}
