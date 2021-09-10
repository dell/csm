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

	"github.com/dell/csm-deployment/model"
)

type driverResponse struct {
	ID                 string `json:"id"`
	StorageArrayTypeID string `json:"storage_array_type_id"`
	Version            string `json:"version"`
} //@name DriverResponse

func newDriverResponse(t *model.DriverType) *driverResponse {
	r := driverResponse{}
	r.ID = fmt.Sprintf("%d", t.ID)
	r.StorageArrayTypeID = fmt.Sprintf("%d", t.StorageArrayTypeID)
	r.Version = t.Version
	return &r
}
