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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

var (
	// Username - Placeholder for username
	Username = ""

	// Password - Placeholder for password
	Password = ""
)

// LoginUser - send call to CSM API for user login
func LoginUser(username, password string) error {
	if password == "" {
		return fmt.Errorf("empty password")
	}
	userLogin := &types.User{
		Token: "",
	}
	err := saveAuthCreds(userLogin)
	if err != nil {
		return err
	}
	Username = username
	Password = password

	userLoginResponse := types.JWTToken
	err = HTTPClient(http.MethodPost, UserLoginURI, nil, &userLoginResponse)
	if err != nil {
		return err
	}
	userLoginToken := &types.User{
		Token: userLoginResponse,
	}
	err = saveAuthCreds(userLoginToken)
	if err != nil {
		return err
	}
	return nil
}

// ChangePassword - send call to CSM API for change password
func ChangePassword(username, currentPassword, newPassword string) error {
	if currentPassword == "" {
		return fmt.Errorf("empty current password")
	}
	if newPassword == "" {
		return fmt.Errorf("empty new password")
	}
	userLogin := &types.User{
		Token: "",
	}
	err := saveAuthCreds(userLogin)
	if err != nil {
		return err
	}
	Username = username
	Password = currentPassword

	err = saveAuthCreds(userLogin)
	if err != nil {
		return err
	}

	err = HTTPClient(http.MethodPatch, fmt.Sprintf(ChangePasswordURI, newPassword), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// saveAuthCreds - save JWT to AUTH_CONFIG_PATH
func saveAuthCreds(userLogin *types.User) error {
	file, _ := json.MarshalIndent(userLogin, "", " ")

	configPath := os.Getenv("AUTH_CONFIG_PATH")
	if configPath == "" {
		return fmt.Errorf("AUTH_CONFIG_PATH not set")
	}
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create auth config directory with error: %v", err)
		}
	}

	err := ioutil.WriteFile(filepath.Join(configPath, "user.json"), file, 0600)
	if err != nil {
		return fmt.Errorf("failed to set user auth creds with error %v", err)
	}
	return nil
}
