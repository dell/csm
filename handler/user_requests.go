// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
)

func newUserResponse(u *model.User) *string {
	r := utils.GenerateJWT(u.Username)
	return &r
}
