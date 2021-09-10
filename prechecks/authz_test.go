// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

// Package  prechecks for application creation test
package prechecks

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, AuthorizationValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"cert.crt"}
			moduleConfig := `
			karaviAuthorizationProxy.rootCertificate.filename=cert.crt
			karaviAuthorizationProxy.proxyAuthzToken.data.access=accessToken 
			karaviAuthorizationProxy.proxyAuthzToken.data.refresh=refreshToken 
			`

			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return true, authorizationValidator, ctrl
		},
		"fail due to missing the set value filename in configuration-files db": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleConfig := `
			karaviAuthorizationProxy.rootCertificate.filename=certNew.crt
			karaviAuthorizationProxy.proxyAuthzToken.data.access=accessToken
			karaviAuthorizationProxy.proxyAuthzToken.data.refresh=refreshToken
			`
			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: []string{"certOld.crt"},
				ModuleConfig:    moduleConfig,
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to setting both filename and data values instead of one": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleConfig := `
			karaviAuthorizationProxy.proxyAuthzToken.filename=token.yaml
			karaviAuthorizationProxy.proxyAuthzToken.data.access=accessToken 
			`
			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: []string{"certOld.crt"},
				ModuleConfig:    moduleConfig,
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to  wrong key=value format": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleConfig := `
			karaviAuthorizationProxy.rootCertificate.filename=certNew.crt=
			karaviAuthorizationProxy.proxyAuthzToken.data.access=accessToken 
			karaviAuthorizationProxy.proxyAuthzToken.data.refresh=refreshToken 
			`
			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: []string{"certOld.crt"},
				ModuleConfig:    moduleConfig,
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to for exit": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: true, Msg: "testing"},
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to missing authorization proxy token values": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: []string{"cert.crt"},
				ModuleConfig:    "",
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to missing authorization proxy certificate values": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleConfig := `
			karaviAuthorizationProxy.proxyAuthzToken.data.access=accessToken 
			karaviAuthorizationProxy.proxyAuthzToken.data.refresh=refreshToken 
			`
			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: []string{"cert.crt"},
				ModuleConfig:    moduleConfig,
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to supplying both proxy authz token filename and data": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"cert.crt"}
			moduleConfig := `
			karaviAuthorizationProxy.rootCertificate.filename=cert.crt
			karaviAuthorizationProxy.proxyAuthzToken.filename=cert.crt
			karaviAuthorizationProxy.proxyAuthzToken.data.access=accessToken 
			karaviAuthorizationProxy.proxyAuthzToken.data.refresh=refreshToken 
			`

			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, authorizationValidator, ctrl
		},
		"fail due to supplying both root certificate filename and data": func(*testing.T) (bool, AuthorizationValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"cert.crt"}
			moduleConfig := `
			karaviAuthorizationProxy.rootCertificate.filename=cert.crt
			karaviAuthorizationProxy.proxyAuthzToken.filename=cert.crt
			karaviAuthorizationProxy.rootCertificate.data=data 
			`

			authorizationValidator := AuthorizationValidator{
				Skip: struct {
					Cond bool
					Msg  string
				}{Cond: false},
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, authorizationValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, authorizationValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, authorizationValidator.Validate())
			} else {
				assert.Error(t, authorizationValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
