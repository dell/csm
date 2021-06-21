package store

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks/cluster_store_interface.go -package=mocks github.com/dell/csm-deployment/store ClusterStoreInterface
type ClusterStoreInterface interface {
	Create(*model.Cluster) error
	GetByID(clusterID uint) (*model.Cluster, error)
	UpdateClusterDetails(u *model.Cluster, details *model.ClusterDetails) (err error)
	Delete(u *model.Cluster) error
	Update(u *model.Cluster) error
	GetAll() ([]model.Cluster, error)
}

type ClusterStore struct {
	db *gorm.DB
}

func NewClusterStore(db *gorm.DB) *ClusterStore {
	return &ClusterStore{
		db: db,
	}
}

func (us *ClusterStore) GetByID(clusterID uint) (*model.Cluster, error) {
	var m model.Cluster
	if err := us.db.Preload(clause.Associations).First(&m, clusterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *ClusterStore) GetAll() ([]model.Cluster, error) {
	var cs []model.Cluster
	if err := us.db.Preload(clause.Associations).Find(&cs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return cs, nil
}

func (us *ClusterStore) Create(u *model.Cluster) (err error) {
	return us.db.Create(u).Error
}

func (us *ClusterStore) UpdateClusterDetails(u *model.Cluster, details *model.ClusterDetails) (err error) {
	return us.db.Model(&u).Association("ClusterDetails").Append(details)
}

func (us *ClusterStore) Update(u *model.Cluster) error {
	return us.db.Save(u).Error
}

func (us *ClusterStore) Delete(u *model.Cluster) error {
	return us.db.Delete(u).Error
}
