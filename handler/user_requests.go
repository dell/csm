package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

type userRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
} //@name UserRegisterRequest

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.Username
	h, err := u.HashPassword(r.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type userLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
} //@name UserLoginRequest

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
	Username string `json:"username"`
	Token    string `json:"token"`
} //@name UserResponse

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.Username = u.Username
	r.Token = utils.GenerateJWT(u.Username)
	return r
}

type userUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
} //@name UserUpdateRequest

func newUserUpdateRequest() *userUpdateRequest {
	return new(userUpdateRequest)
}

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
	if r.Password != u.Password {
		h, err := u.HashPassword(r.Password)
		if err != nil {
			return err
		}
		u.Password = h
	}
	return nil
}
