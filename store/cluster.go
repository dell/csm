package store

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

// ClusterStoreInterface is used to define the interface for persisting Clusters
//go:generate mockgen -destination=mocks/cluster_store_interface.go -package=mocks github.com/dell/csm-deployment/store ClusterStoreInterface
type ClusterStoreInterface interface {
	Create(*model.Cluster) error
	GetByID(clusterID uint) (*model.Cluster, error)
	UpdateClusterDetails(u *model.Cluster, details *model.ClusterDetails) (err error)
	Delete(u *model.Cluster) error
	Update(u *model.Cluster) error
	GetAll() ([]model.Cluster, error)
	GetAllByName(string) ([]model.Cluster, error)
}

// ClusterStore is used to operate on the Clusters persistent store
type ClusterStore struct {
	db *gorm.DB
}

// NewClusterStore creates a new ClusterStore
func NewClusterStore(db *gorm.DB) *ClusterStore {
	return &ClusterStore{
		db: db,
	}
}

// GetAllByName will return all clusters with a matching name
func (us *ClusterStore) GetAllByName(name string) ([]model.Cluster, error) {
	var clusters []model.Cluster
	if err := us.db.
		Preload(clause.Associations).
		Where(&model.Cluster{ClusterName: name}).
		First(&clusters).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return clusters, nil
}

// GetByID returns a cluster with the given ID
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

// GetAll returns all clusters
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

// Create will persist a new cluster in the database
func (us *ClusterStore) Create(u *model.Cluster) (err error) {
	return us.db.Create(u).Error
}

// UpdateClusterDetails will update the ClusterDetails association in the database
func (us *ClusterStore) UpdateClusterDetails(u *model.Cluster, details *model.ClusterDetails) (err error) {
	return us.db.Model(&u).Association("ClusterDetails").Append(details)
}

// Update will update an existing cluster in the database
func (us *ClusterStore) Update(u *model.Cluster) error {
	return us.db.Save(u).Error
}

// Delete will delete an existing cluster from the database
func (us *ClusterStore) Delete(u *model.Cluster) error {
	return us.db.Delete(u).Error
}
