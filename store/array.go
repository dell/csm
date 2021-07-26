package store

import (
	"errors"
	"fmt"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StorageArrayStoreInterface  is used to define the interface for persisting Array type
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

// StorageArrayStore is used to operate on the Storage Array persistent store
type StorageArrayStore struct {
	db *gorm.DB
}

// NewStorageArrayStore creates new StorageArrayStore
func NewStorageArrayStore(db *gorm.DB) *StorageArrayStore {
	return &StorageArrayStore{
		db: db,
	}
}

// GetByID -  returns array by Id
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

// GetAllByID - returns arrays by Id
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

// GetAll will return all storage arrays
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

// GetTypeByTypeName - get array type by name
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

// Create - Method to add new Array
func (sas *StorageArrayStore) Create(a *model.StorageArray) error {
	return sas.db.Create(a).Error
}

// Update - Method to update array info
func (sas *StorageArrayStore) Update(a *model.StorageArray) (err error) {
	return sas.db.Save(a).Error
}

// Delete - Method to delete array
func (sas *StorageArrayStore) Delete(a *model.StorageArray) error {
	return sas.db.Delete(a).Error
}
