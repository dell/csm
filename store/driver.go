package store

import (
	"errors"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination=mocks/driver_type_store_interface.go -package=mocks github.com/dell/csm-deployment/store DriverTypeStoreInterface
type DriverTypeStoreInterface interface {
	GetByID(uint) (*model.DriverType, error)
	GetAll() ([]model.DriverType, error)
}

type DriverTypeStore struct {
	db *gorm.DB
}

func NewDriverTypeStore(db *gorm.DB) *DriverTypeStore {
	return &DriverTypeStore{
		db: db,
	}
}

func (sas *DriverTypeStore) GetByID(id uint) (*model.DriverType, error) {
	var sa model.DriverType
	if err := sas.db.Preload(clause.Associations).First(&sa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sa, nil
}

func (sas *DriverTypeStore) GetAll() ([]model.DriverType, error) {
	var sa []model.DriverType
	if err := sas.db.Preload(clause.Associations).Find(&sa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}
