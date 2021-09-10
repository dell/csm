// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
)

type applicationCreateRequest struct {
	Name                string   `json:"name" validate:"required"`
	ClusterID           string   `json:"cluster_id" validate:"required"`
	DriverTypeID        string   `json:"driver_type_id" validate:"required"`
	ModuleTypes         []string `json:"module_types"`
	StorageArrays       []string `json:"storage_arrays"`
	DriverConfiguration []string `json:"driver_configuration"`
	ModuleConfiguration []string `json:"module_configuration"`
} //@name ApplicationCreateRequest

type applicationResponse struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	ClusterID           string   `json:"cluster_id"`
	DriverTypeID        string   `json:"driver_type_id"`
	ModuleTypes         []string `json:"module_types"`
	StorageArrays       []string `json:"storage_arrays"`
	DriverConfiguration []string `json:"driver_configuration"`
	ModuleConfiguration []string `json:"module_configuration"`
	ApplicationOutput   string   `json:"application_output"`
} //@name ApplicationResponse

func newApplicationResponse(a *model.Application) *applicationResponse {
	r := new(applicationResponse)
	r.ID = fmt.Sprintf("%d", a.ID)
	r.Name = a.Name
	r.ClusterID = fmt.Sprintf("%d", a.ClusterID)
	r.DriverTypeID = fmt.Sprintf("%d", a.DriverTypeID)
	for _, v := range a.ModuleTypes {
		r.ModuleTypes = append(r.ModuleTypes, fmt.Sprintf("%d", v.ID))
	}
	for _, v := range a.StorageArrays {
		r.StorageArrays = append(r.StorageArrays, fmt.Sprintf("%d", v.ID))
	}
	r.ApplicationOutput = a.ApplicationOutput
	r.DriverConfiguration = strings.Split(a.DriverConfiguration, " ")
	r.ModuleConfiguration = strings.Split(a.ModuleConfiguration, " ")
	return r
}

func (r *applicationCreateRequest) bind(c echo.Context, application *model.Application, moduleTypeStore store.ModuleTypeStoreInterface) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	application.Name = r.Name
	clusterID, err := strconv.Atoi(r.ClusterID)
	if err != nil {
		return err
	}
	application.ClusterID = uint(clusterID)

	driverTypeID, err := strconv.Atoi(r.DriverTypeID)
	if err != nil {
		return err
	}

	application.DriverTypeID = uint(driverTypeID)
	application.ModuleTypes = make([]model.ModuleType, 0)
	for _, moduleTypeID := range r.ModuleTypes {
		moduleTypeID, err := strconv.Atoi(moduleTypeID)
		if err != nil {
			return err
		}
		moduleType, err := moduleTypeStore.GetByID(uint(moduleTypeID))
		if err != nil {
			return err
		}
		application.ModuleTypes = append(application.ModuleTypes, *moduleType)
	}
	application.DriverConfiguration = strings.Join(r.DriverConfiguration, " ")
	application.ModuleConfiguration = strings.Join(r.ModuleConfiguration, " ")

	return nil
}
