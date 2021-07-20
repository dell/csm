package store

import (
	"errors"
	"fmt"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination=mocks/storage_array_store_interface.go -package=mocks github.com/dell/csm-deployment/store StorageArrayStoreInterface
type StorageArrayStoreInterface interface {
	GetByID(uint) (*model.StorageArray, error)
	GetAllByID(...uint) ([]model.StorageArray, error)
	GetTypeByTypeName(string) (*model.StorageArrayType, error)
	Create(*model.StorageArray) error
	GetAll() ([]model.StorageArray, error)
	GetAllByUniqueID(string) ([]model.StorageArray, error)
	GetAllByStorageType(string) ([]model.StorageArray, error)
	Delete(*model.StorageArray) error
	Update(*model.StorageArray) error
}

type StorageArrayStore struct {
	db *gorm.DB
}

func NewStorageArrayStore(db *gorm.DB) *StorageArrayStore {
	return &StorageArrayStore{
		db: db,
	}
}

func (sas *StorageArrayStore) GetByID(id uint) (*model.StorageArray, error) {
	var sa model.StorageArray
	if err := sas.db.Preload(clause.Associations).First(&sa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sa, nil
}

func (sas *StorageArrayStore) GetAllByID(v ...uint) ([]model.StorageArray, error) {
	var sa []model.StorageArray
	if len(v) > 0 {
		if err := sas.db.Preload(clause.Associations).Find(&sa, v).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}
	}
	return sa, nil
}

func (sas *StorageArrayStore) GetAll() ([]model.StorageArray, error) {
	var sa []model.StorageArray
	if err := sas.db.Preload(clause.Associations).Find(&sa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}

// GetAllByUniqueID will return all storage arrays with matching unique id
func (sas *StorageArrayStore) GetAllByUniqueID(uniqueID string) ([]model.StorageArray, error) {
	var storageArrays []model.StorageArray
	if err := sas.db.
		Preload(clause.Associations).
		Where(&model.StorageArray{UniqueID: uniqueID}).
		Find(&storageArrays).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return storageArrays, nil
}

// GetAllByStorageType will return all storage arrays with matching storage type
func (sas *StorageArrayStore) GetAllByStorageType(storageTypeName string) ([]model.StorageArray, error) {
	var storageArrays []model.StorageArray
	storageType, err := sas.GetTypeByTypeName(storageTypeName)
	if err != nil {
		return nil, err
	}
	if storageType == nil {
		return nil, fmt.Errorf("unable to find storage type with name %s", storageTypeName)
	}
	if err := sas.db.
		Preload(clause.Associations).
		Where(&model.StorageArray{StorageArrayTypeID: storageType.ID}).
		Find(&storageArrays).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return storageArrays, nil
}

func (sas *StorageArrayStore) GetTypeByTypeName(typeName string) (*model.StorageArrayType, error) {
	var sat model.StorageArrayType
	if err := sas.db.
		Preload(clause.Associations).
		Where(&model.StorageArrayType{Name: typeName}).
		First(&sat).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sat, nil
}

func (sas *StorageArrayStore) Create(a *model.StorageArray) error {
	return sas.db.Create(a).Error
}

func (sas *StorageArrayStore) Update(a *model.StorageArray) (err error) {
	return sas.db.Save(a).Error
}

func (sas *StorageArrayStore) Delete(a *model.StorageArray) error {
	return sas.db.Delete(a).Error
}
