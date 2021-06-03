package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type userLoginRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type userResponse struct {
	User struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Token = utils.GenerateJWT(u.Username)
	return r
}

type userUpdateRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

func newUserUpdateRequest() *userUpdateRequest {
	return new(userUpdateRequest)
}

func (r *userUpdateRequest) populate(u *model.User) {
	r.User.Username = u.Username
	r.User.Password = u.Password
}

func (r *userUpdateRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	if r.User.Password != u.Password {
		h, err := u.HashPassword(r.User.Password)
		if err != nil {
			return err
		}
		u.Password = h
	}
	return nil
}
