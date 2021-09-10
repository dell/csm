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

type moduleResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Standalone bool   `json:"standalone"`
} //@name ModuleResponse

func newModuleResponse(t *model.ModuleType) *moduleResponse {
	r := moduleResponse{}
	r.ID = fmt.Sprintf("%d", t.ID)
	r.Name = t.Name
	r.Version = t.Version
	r.Standalone = t.Standalone
	return &r
}
