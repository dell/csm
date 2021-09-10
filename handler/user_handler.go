// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
)

// UserHandler constains the store interface for User
type UserHandler struct {
	userStore store.UserStoreInterface
}

// New creates an handler for a User
func New(us store.UserStoreInterface) *UserHandler {
	return &UserHandler{
		userStore: us,
	}
}

// Register rosters all the API endpoints for users
func (h *UserHandler) Register(api *echo.Group) {
	adminUsers := api.Group("/users")
	adminUsers.POST("/login", h.Login)
	adminUsers.PATCH("/change-password", h.ChangePasword)
}
