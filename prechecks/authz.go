// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package prechecks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dell/csm-deployment/utils"
	"github.com/pkg/errors"
)

// AuthorizationValidator validates the required precheck are met for deploying authorization
type AuthorizationValidator struct {
	ModuleConfig    string
	ConfigFileNames []string
	Skip            struct {
		Cond bool
		Msg  string
	}
}

// Validate will check that all the expected cert-manager components exist and are running
func (k AuthorizationValidator) Validate() error {
	if k.Skip.Cond {
		return errors.New(k.Skip.Msg)
	}

	if !(strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.proxyAuthzToken.filename") ||
		(strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.proxyAuthzToken.data.access") &&
			strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.proxyAuthzToken.data.refresh"))) {
		return errors.New("missing authorization proxy token")
	}

	if !(strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.rootCertificate.filename") ||
		strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.rootCertificate.data")) {
		return errors.New("missing authorization proxy certificate")
	}

	if strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.proxyAuthzToken.filename") &&
		(strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.proxyAuthzToken.data.access") ||
			strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.proxyAuthzToken.data.refresh")) {
		return errors.New("both filename and data are set for proxyAuthzToken, set either but not both")
	}

	if strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.rootCertificate.filename") &&
		strings.Contains(k.ModuleConfig, "karaviAuthorizationProxy.rootCertificate.data") {
		return errors.New("both filename and data are set for rootCertificate, set either but not both")
	}

	re := regexp.MustCompile(`[^\s]+`)
	for _, v := range re.FindAllString(k.ModuleConfig, -1) {
		if strings.Contains(v, ".filename") {
			configFile := strings.Split(v, "=")
			if len(configFile) != 2 {
				return errors.New("invalid ytt value format. It should be key1=value1")
			}
			if !utils.Find(k.ConfigFileNames, configFile[1]) {
				return fmt.Errorf("the filename %s for the module configuration value %s does not exist in configuration-files", configFile[1], configFile[0])
			}
		}
	}

	return nil
}
