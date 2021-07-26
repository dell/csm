package store

import (
	"errors"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination=mocks/module_type_store_interface.go -package=mocks github.com/dell/csm-deployment/store ModuleTypeStoreInterface
// ModuleTypeStoreInterface is used to define the interface for persisting Modules
type ModuleTypeStoreInterface interface {
	GetByID(uint) (*model.ModuleType, error)
	GetAll() ([]model.ModuleType, error)
	GetAllByID(...uint) ([]model.ModuleType, error)
}

// ModuleTypeStore - Placeholder for Module Type Store
type ModuleTypeStore struct {
	db *gorm.DB
}

// NewModuleTypeStore returns an instance of ModuleTypeStore in db
func NewModuleTypeStore(db *gorm.DB) *ModuleTypeStore {
	return &ModuleTypeStore{
		db: db,
	}
}

// GetByID returns an instance of ModuleTypeStore  that matches id in db
func (ms *ModuleTypeStore) GetByID(id uint) (*model.ModuleType, error) {
	var mt model.ModuleType
	if err := ms.db.Preload(clause.Associations).First(&mt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

// GetAll returns all instances of ModuleTypeStore in db
func (ms *ModuleTypeStore) GetAll() ([]model.ModuleType, error) {
	var sa []model.ModuleType
	if err := ms.db.Preload(clause.Associations).Find(&sa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}

// GetAllByID returns instances of ModuleTypeStore in db that match all passed in IDs
func (ms *ModuleTypeStore) GetAllByID(v ...uint) ([]model.ModuleType, error) {
	var mt []model.ModuleType
	if len(v) > 0 {
		if err := ms.db.Preload(clause.Associations).Find(&mt, v).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}
	}
	return mt, nil
}
