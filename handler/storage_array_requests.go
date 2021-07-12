package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

type storageArrayCreateRequest struct {
	StorageArrayType   string `json:"storage_array_type" validate:"required"`
	UniqueID           string `json:"unique_id" validate:"required"`
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password" validate:"required"`
	ManagementEndpoint string `json:"management_endpoint" validate:"required"`
} //@name StorageArrayCreateRequest

type storageArrayUpdateRequest struct {
	StorageArrayType   string `json:"storage_array_type" validate:"required"`
	UniqueID           string `json:"unique_id" validate:"required"`
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password" validate:"required"`
	ManagementEndpoint string `json:"management_endpoint" validate:"required"`
} //@name StorageArrayUpdateRequest

type storageArrayResponse struct {
	ID                 uint   `json:"id"`
	StorageArrayTypeID uint   `json:"storage_array_type_id"`
	UniqueID           string `json:"unique_id"`
	Username           string `json:"username"`
	ManagementEndpoint string `json:"management_endpoint"`
} //@name StorageArrayResponse

func newStorageArrayResponse(arr *model.StorageArray) *storageArrayResponse {
	r := new(storageArrayResponse)
	r.ID = arr.ID
	r.UniqueID = arr.UniqueID
	r.Username = arr.Username
	r.StorageArrayTypeID = arr.StorageArrayTypeID
	r.ManagementEndpoint = arr.ManagementEndpoint
	return r
}

func (r *storageArrayCreateRequest) bind(c echo.Context, array *model.StorageArray) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	array.UniqueID = r.UniqueID
	array.Username = r.Username
	array.ManagementEndpoint = r.ManagementEndpoint

	// TODO: it better to store hash, but we will need password for secret creation
	encrypted, err := utils.EncryptPassword([]byte(r.Password))
	if err != nil {
		return err
	}
	array.Password = string(encrypted)
	return nil
}

func (r *storageArrayUpdateRequest) bind(c echo.Context, array *model.StorageArray) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	array.UniqueID = r.UniqueID
	array.Username = r.Username
	array.ManagementEndpoint = r.ManagementEndpoint

	// TODO: it better to store hash, but we will need password for secret creation
	encrypted, err := utils.EncryptPassword([]byte(r.Password))
	if err != nil {
		return err
	}
	array.Password = string(encrypted)
	return nil
}
