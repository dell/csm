package store

import (
	"errors"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination=mocks/module_type_store_interface.go -package=mocks github.com/dell/csm-deployment/store ModuleStoreInterface
type ModuleStoreInterface interface {
	GetByID(uint) (*model.ModuleType, error)
	GetAll() ([]model.ModuleType, error)
	GetAllByID(...uint) ([]model.ModuleType, error)
}

type ModuleStore struct {
	db *gorm.DB
}

// NewModuleStore returns an instance of ModuleStore in db
func NewModuleStore(db *gorm.DB) *ModuleStore {
	return &ModuleStore{
		db: db,
	}
}

// GetByID returns an instance of ModuleStore  that matches id in db
func (ms *ModuleStore) GetByID(id uint) (*model.ModuleType, error) {
	var mt model.ModuleType
	if err := ms.db.Preload(clause.Associations).First(&mt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

// GetAll returns all instances of ModuleStore in db
func (sas *ModuleStore) GetAll() ([]model.ModuleType, error) {
	var sa []model.ModuleType
	if err := sas.db.Preload(clause.Associations).Find(&sa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}

// GetAllByID returns instances of ModuleStore in db that match all passed in IDs
func (ms *ModuleStore) GetAllByID(v ...uint) ([]model.ModuleType, error) {
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
