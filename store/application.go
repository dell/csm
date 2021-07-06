package store

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

// ApplicationStoreInterface is used to define the interface for persisting Application
//go:generate mockgen -destination=mocks/application_store_interface.go -package=mocks github.com/dell/csm-deployment/store ApplicationStoreInterface
type ApplicationStoreInterface interface {
	Create(*model.Application) error
	GetByID(string) (*model.Application, error)
	Update(*model.Application) error
	Delete(a *model.Application) error
}

type ApplicationStore struct {
	db *gorm.DB
}

func NewApplicationStore(db *gorm.DB) *ApplicationStore {
	return &ApplicationStore{
		db: db,
	}
}

func (as *ApplicationStore) GetByID(id string) (*model.Application, error) {
	var m model.Application
	if err := as.db.Preload(clause.Associations).Preload("StorageArrays.StorageArrayType").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (as *ApplicationStore) Create(a *model.Application) (err error) {
	return as.db.Save(a).Error
}

func (as *ApplicationStore) Update(a *model.Application) (err error) {
	return as.db.Save(a).Error
}

func (as *ApplicationStore) Delete(a *model.Application) error {
	return as.db.Delete(a).Error
}
