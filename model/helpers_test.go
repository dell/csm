// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package model_test

import (
	"github.com/dell/csm-deployment/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HashPassword(t *testing.T) {
	tests := map[string]func(t *testing.T) (model.User, *gomock.Controller){
		"success - HashPassword": func(*testing.T) (model.User, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			user := model.User{Username: "user", Password: "password"}
			return user, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			user, ctrl := tc(t)
			hashedPassword, err := user.HashPassword(user.Password)
			assert.NoError(t, err)
			assert.NotNil(t, hashedPassword, "Hashed Password should not be nil")
			user.Password = hashedPassword
			assert.True(t, user.CheckPassword("password"))
			ctrl.Finish()
		})
	}
}

func Test_HashPasswordFailure(t *testing.T) {
	tests := map[string]func(t *testing.T) (model.User, *gomock.Controller){
		"failure - HashPassword": func(*testing.T) (model.User, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			user := model.User{Username: "user", Password: ""}
			return user, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			user, ctrl := tc(t)
			_, err := user.HashPassword(user.Password)
			assert.Error(t, err)
			ctrl.Finish()
		})
	}
}
