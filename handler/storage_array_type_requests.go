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

type storageArrayTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} //@name StorageArrayTypeResponse

func newStorageArrayTypeResponse(t *model.StorageArrayType) *storageArrayTypeResponse {
	r := storageArrayTypeResponse{}
	r.ID = fmt.Sprintf("%d", t.ID)
	r.Name = t.Name
	return &r
}
