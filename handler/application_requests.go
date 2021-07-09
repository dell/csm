package handler

import (
	"strings"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
)

type applicationCreateRequest struct {
	Application struct {
		Name                string   `json:"name" validate:"required"`
		ClusterID           uint     `json:"cluster_id"`
		DriverTypeID        uint     `json:"driver_type_id"`
		ModuleTypes         []uint   `json:"module_types"`
		StorageArrays       []uint   `json:"storage_arrays"`
		DriverConfiguration []string `json:"driver_configuration"`
		ModuleConfiguration []string `json:"module_configuration"`
	} `json:"application"`
} //@name ApplicationCreateRequest

type applicationResponse struct {
	Application struct {
		ID                  uint     `json:"id"`
		Name                string   `json:"name"`
		ClusterID           uint     `json:"cluster_id"`
		DriverTypeID        uint     `json:"driver_type_id"`
		ModuleTypes         []uint   `json:"module_types"`
		StorageArrays       []uint   `json:"storage_arrays"`
		DriverConfiguration []string `json:"driver_configuration"`
		ModuleConfiguration []string `json:"module_configuration"`
		ApplicationOutput   string   `json:"application_output"`
	} `json:"application"`
} //@name ApplicationResponse

func newApplicationResponse(a *model.Application) *applicationResponse {
	r := new(applicationResponse)
	r.Application.ID = a.ID
	r.Application.Name = a.Name
	r.Application.ClusterID = a.ClusterID
	r.Application.DriverTypeID = a.DriverTypeID
	for _, v := range a.ModuleTypes {
		r.Application.ModuleTypes = append(r.Application.ModuleTypes, v.ID)
	}
	for _, v := range a.StorageArrays {
		r.Application.StorageArrays = append(r.Application.StorageArrays, v.ID)
	}
	r.Application.ApplicationOutput = a.ApplicationOutput
	r.Application.DriverConfiguration = strings.Split(a.DriverConfiguration, " ")
	r.Application.ModuleConfiguration = strings.Split(a.ModuleConfiguration, " ")
	return r
}

func (r *applicationCreateRequest) bind(c echo.Context, application *model.Application, moduleTypeStore store.ModuleStoreInterface) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	application.Name = r.Application.Name
	application.ClusterID = r.Application.ClusterID
	application.DriverTypeID = r.Application.DriverTypeID
	application.ModuleTypes = make([]model.ModuleType, 0)
	for _, moduleTypeID := range r.Application.ModuleTypes {
		moduleType, err := moduleTypeStore.GetByID(moduleTypeID)
		if err != nil {
			return err
		}
		application.ModuleTypes = append(application.ModuleTypes, *moduleType)
	}
	application.DriverConfiguration = strings.Join(r.Application.DriverConfiguration, " ")
	application.ModuleConfiguration = strings.Join(r.Application.ModuleConfiguration, " ")

	return nil
}
