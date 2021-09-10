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

type taskResponse struct {
	ID              string                       `json:"id"`
	Status          string                       `json:"status"`
	ApplicationName string                       `json:"application_name"`
	Logs            string                       `json:"logs"`
	Links           map[string]map[string]string `json:"_links"`
} //@name TaskResponse

func newTaskResponse(t *model.Task) *taskResponse {
	r := taskResponse{}
	r.ID = fmt.Sprintf("%d", t.ID)
	r.Status = t.Status
	r.ApplicationName = t.Application.Name
	r.Logs = string(t.Logs)
	return &r
}

func newTaskResponseWithLinks(t *model.Task, links map[string]map[string]string) *taskResponse {
	r := newTaskResponse(t)
	r.Links = links
	return r
}
