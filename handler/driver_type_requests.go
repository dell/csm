package handler

import "github.com/dell/csm-deployment/model"

type driverResponse struct {
	ID                 uint   `json:"id"`
	StorageArrayTypeID uint   `json:"storage_array_type_ID"`
	Version            string `json:"version"`
} //@name DriverResponse

func newDriverResponse(t *model.DriverType) *driverResponse {
	r := driverResponse{}
	r.ID = t.ID
	r.StorageArrayTypeID = t.StorageArrayTypeID
	r.Version = t.Version
	return &r
}
