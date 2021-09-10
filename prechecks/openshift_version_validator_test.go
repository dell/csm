// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package prechecks

import (
	"errors"
	"testing"

	"github.com/dell/csm-deployment/prechecks/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_OpenshiftVersionValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.7", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "4.8",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return true, versionValidator, ctrl
		},
		"success - at minimum version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.6", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "4.7",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return true, versionValidator, ctrl
		},
		"success - at maximum version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.7", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "4.7",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return true, versionValidator, ctrl
		},
		"success - skip k8s": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)

			versionValidator := OpenshiftVersionValidator{
				K8sClient: versionInterface,
				Logger:    echo.New().Logger,
			}

			return true, versionValidator, ctrl
		},
		"error - below minimum version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.5", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "4.7",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
		"error - above maximum minimum version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.8", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "4.7",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
		"error - checking openshift": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, errors.New("error"))

			versionValidator := OpenshiftVersionValidator{
				K8sClient: versionInterface,
				Logger:    echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
		"error - getting version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("", errors.New("error"))

			versionValidator := OpenshiftVersionValidator{
				K8sClient: versionInterface,
				Logger:    echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
		"error - invalid min version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.6", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "invalid-version",
				MaximumVersion: "4.7",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
		"error - invalid max version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("4.6", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "invalid-version",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
		"error - invalid version": func(*testing.T) (bool, OpenshiftVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("invalid-version", nil)

			versionValidator := OpenshiftVersionValidator{
				MinimumVersion: "4.6",
				MaximumVersion: "4.7",
				K8sClient:      versionInterface,
				Logger:         echo.New().Logger,
			}

			return false, versionValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, versionValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, versionValidator.Validate())
			} else {
				assert.Error(t, versionValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
