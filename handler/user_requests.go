package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

const (
	// DefaultUsername is the default K8s admin username
	DefaultUsername = "admin"
	// DefaultPassword is the default K8s admin password
	DefaultPassword = "password"
)

type userResponse struct {
	Token    string   `json:"token"`
	Messages []string `json:"mesages"`
} //@name UserResponse

func newUserResponse(u *model.User) *userResponse {
	r := userResponse{
		Token:    utils.GenerateJWT(u.Username),
		Messages: []string{},
	}

	if u.Username == DefaultUsername {
		r.Messages = append(r.Messages, "username still matches the default value. Consider changing username(HIGHLY RECOMMENDED).")
	}
	if u.Password == DefaultPassword {
		r.Messages = append(r.Messages, "password still matches the default value. Consider changing password(HIGHLY RECOMMENDED).")
	}

	return &r
}

type userUpdateRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
} //@name UserUpdateRequest

func (r *userUpdateRequest) populate(u *model.User) {
	r.Username = u.Username
	r.Password = u.Password
}

func (r *userUpdateRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.Username
	u.Password = r.Password
	return nil
}
