package store

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

type ClusterStoreInterface interface {
	Create(*model.Cluster) error
	GetByClusterID(clusterID uint) (*model.Cluster, error)
	UpdateClusterDetails(u *model.Cluster, details *model.ClusterDetails) (err error)
}

type ClusterStore struct {
	db *gorm.DB
}

func NewClusterStore(db *gorm.DB) *ClusterStore {
	return &ClusterStore{
		db: db,
	}
}

func (us *ClusterStore) GetByClusterID(clusterID uint) (*model.Cluster, error) {
	var m model.Cluster
	if err := us.db.Preload(clause.Associations).First(&m, clusterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *ClusterStore) Create(u *model.Cluster) (err error) {
	return us.db.Create(u).Error
}

func (us *ClusterStore) UpdateClusterDetails(u *model.Cluster, details *model.ClusterDetails) (err error) {
	return us.db.Model(&u).Association("ClusterDetails").Append(details)
}
