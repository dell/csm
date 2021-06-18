package store

import (
	"errors"
	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModuleStoreInterface interface {
	GetByID(uint) (*model.ModuleType, error)
	GetAllByID(...uint) ([]model.ModuleType, error)
}

type ModuleStore struct {
	db *gorm.DB
}

func NewModuleStore(db *gorm.DB) *ModuleStore {
	return &ModuleStore{
		db: db,
	}
}

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
