package store

import (
	"errors"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

// TaskStoreInterface is used to define the interface for persisting Task
//go:generate mockgen -destination=mocks/task_store_interface.go -package=mocks github.com/dell/csm-deployment/store TaskStoreInterface
type TaskStoreInterface interface {
	Create(*model.Task) error
	GetByID(string) (*model.Task, error)
	Update(*model.Task) error
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
	if err := ts.db.First(&t, id).Error; err != nil {
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
