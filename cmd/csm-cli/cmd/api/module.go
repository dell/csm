// Package api for API services
// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

// GetModuleTypeID - returns ID of input module type
func GetModuleTypeID(moduleType string) (string, error) {
	getModulesResponse, err := GetModuleTypes()
	if err != nil {
		return "", fmt.Errorf("failed to get supported modules with error: %v", err)
	}
	// moduleType is a combination of name and version
	// split moduleType to get corresponding id
	moduleTypeSplit := strings.Split(moduleType, ":")
	name := ""
	if len(moduleTypeSplit) > 0 {
		name = moduleTypeSplit[0]
	}
	version := ""
	if len(moduleTypeSplit) == 2 {
		version = strings.Trim(moduleTypeSplit[1], "v")
	}
	if len(moduleTypeSplit) <= 2 {
		for _, module := range getModulesResponse {
			if module.Name == name && module.Version == version {
				return module.ID, nil
			}
		}
	}
	return "", fmt.Errorf("invalid module type: %s", moduleType)
}

// GetModuleTypes - returns supported module types by CSM API
func GetModuleTypes() ([]types.ModuleTypeResponse, error) {
	getModuleTypeResponse := []types.ModuleTypeResponse{}
	err := HTTPClient(http.MethodGet, GetModuleTypeURI, nil, &getModuleTypeResponse)
	if err != nil {
		return nil, err
	}
	return getModuleTypeResponse, nil
}
