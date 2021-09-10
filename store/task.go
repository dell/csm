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
	"gorm.io/gorm/clause"
)

// TaskStoreInterface is used to define the interface for persisting Task
//go:generate mockgen -destination=mocks/task_store_interface.go -package=mocks github.com/dell/csm-deployment/store TaskStoreInterface
type TaskStoreInterface interface {
	Create(*model.Task) error
	GetByID(string) (*model.Task, error)
	Update(*model.Task) error
	GetAll() ([]model.Task, error)
	GetAllByApplication(uint) ([]model.Task, error)
}

// TaskStore - Placeholder for Task store
type TaskStore struct {
	db *gorm.DB
}

// NewTaskStore - returns an instance of TaskStore in db
func NewTaskStore(db *gorm.DB) *TaskStore {
	return &TaskStore{
		db: db,
	}
}

// GetByID - Method to get Task by Id
func (ts *TaskStore) GetByID(id string) (*model.Task, error) {
	var t model.Task
	if err := ts.db.Preload(clause.Associations).First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Create - Method to create Task
func (ts *TaskStore) Create(t *model.Task) error {
	return ts.db.Create(t).Error
}

// Update - Method to update Task
func (ts *TaskStore) Update(t *model.Task) error {
	return ts.db.Model(t).Updates(t).Error
}

// GetAllByApplication will return all tasks for a given application
func (ts *TaskStore) GetAllByApplication(id uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := ts.db.
		Preload(clause.Associations).
		Where(&model.Task{ApplicationID: id}).
		Find(&tasks).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tasks, nil
}

// GetAll returns all tasks
func (ts *TaskStore) GetAll() ([]model.Task, error) {
	var tasks []model.Task
	if err := ts.db.Preload(clause.Associations).Find(&tasks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tasks, nil
}
