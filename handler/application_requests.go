package handler

import (
	"strings"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
)

type applicationCreateRequest struct {
	Name                string   `json:"name" validate:"required"`
	ClusterID           uint     `json:"cluster_id"`
	DriverTypeID        uint     `json:"driver_type_id"`
	ModuleTypes         []uint   `json:"module_types"`
	StorageArrays       []uint   `json:"storage_arrays"`
	DriverConfiguration []string `json:"driver_configuration"`
	ModuleConfiguration []string `json:"module_configuration"`
} //@name ApplicationCreateRequest

type applicationResponse struct {
	ID                  uint     `json:"id"`
	Name                string   `json:"name"`
	ClusterID           uint     `json:"cluster_id"`
	DriverTypeID        uint     `json:"driver_type_id"`
	ModuleTypes         []uint   `json:"module_types"`
	StorageArrays       []uint   `json:"storage_arrays"`
	DriverConfiguration []string `json:"driver_configuration"`
	ModuleConfiguration []string `json:"module_configuration"`
	ApplicationOutput   string   `json:"application_output"`
} //@name ApplicationResponse

func newApplicationResponse(a *model.Application) *applicationResponse {
	r := new(applicationResponse)
	r.ID = a.ID
	r.Name = a.Name
	r.ClusterID = a.ClusterID
	r.DriverTypeID = a.DriverTypeID
	for _, v := range a.ModuleTypes {
		r.ModuleTypes = append(r.ModuleTypes, v.ID)
	}
	for _, v := range a.StorageArrays {
		r.StorageArrays = append(r.StorageArrays, v.ID)
	}
	r.ApplicationOutput = a.ApplicationOutput
	r.DriverConfiguration = strings.Split(a.DriverConfiguration, " ")
	r.ModuleConfiguration = strings.Split(a.ModuleConfiguration, " ")
	return r
}

func (r *applicationCreateRequest) bind(c echo.Context, application *model.Application, moduleTypeStore store.ModuleStoreInterface) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	application.Name = r.Name
	application.ClusterID = r.ClusterID
	application.DriverTypeID = r.DriverTypeID
	application.ModuleTypes = make([]model.ModuleType, 0)
	for _, moduleTypeID := range r.ModuleTypes {
		moduleType, err := moduleTypeStore.GetByID(moduleTypeID)
		if err != nil {
			return err
		}
		application.ModuleTypes = append(application.ModuleTypes, *moduleType)
	}
	application.DriverConfiguration = strings.Join(r.DriverConfiguration, " ")
	application.ModuleConfiguration = strings.Join(r.ModuleConfiguration, " ")

	return nil
}
