package store

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

// ApplicationStateChangeStoreInterface is used to define the interface for persisting Application State Change
//go:generate mockgen -destination=mocks/application_state_change_store_interface.go -package=mocks github.com/dell/csm-deployment/store ApplicationStateChangeStoreInterface
type ApplicationStateChangeStoreInterface interface {
	Create(*model.ApplicationStateChange) error
	GetByApplicationID(uint) (*model.ApplicationStateChange, error)
	GetById(id uint) (*model.ApplicationStateChange, error)
	Delete(a *model.ApplicationStateChange) error
}

// ApplicationStateChangeStore - Placeholder for Application state change store
type ApplicationStateChangeStore struct {
	db *gorm.DB
}

// NewApplicationStateChangeStore - returns an instance of ApplicationStateChangeStore in db
func NewApplicationStateChangeStore(db *gorm.DB) *ApplicationStateChangeStore {
	return &ApplicationStateChangeStore{
		db: db,
	}
}

// GetById - returns Application by Id
func (as *ApplicationStateChangeStore) GetById(id uint) (*model.ApplicationStateChange, error) {
	var m model.ApplicationStateChange
	if err := as.db.Preload(clause.Associations).Preload("StorageArrays.StorageArrayType").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// GetByApplicationID - returns Application by Id
func (as *ApplicationStateChangeStore) GetByApplicationID(id uint) (*model.ApplicationStateChange, error) {
	var m model.ApplicationStateChange
	if err := as.db.Where(&model.ApplicationStateChange{ApplicationID: id}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// Create new Application State Change
func (as *ApplicationStateChangeStore) Create(a *model.ApplicationStateChange) (err error) {
	return as.db.Save(a).Error
}

// Delete Application State Change
func (as *ApplicationStateChangeStore) Delete(a *model.ApplicationStateChange) error {
	return as.db.Delete(a).Error
}
