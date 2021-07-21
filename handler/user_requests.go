package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
)

func newUserResponse(u *model.User) *string {
	r := utils.GenerateJWT(u.Username)
	return &r
}
