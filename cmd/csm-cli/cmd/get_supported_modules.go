// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package cmd

import (
	"fmt"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
)

// getSupportedModulesCmd represents the supported-modules command
var getSupportedModulesCmd = &cobra.Command{
	Use:   "supported-modules",
	Short: "get list of supported modules",
	Long:  `get list of supported modules`,
	Run: func(cmd *cobra.Command, args []string) {
		log := utils.GetLogger()
		supportedModules, err := api.GetModuleTypes()
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("get supported module failed")
		} else if len(supportedModules) == 0 {
			fmt.Println("no supported modules")
		} else {
			for _, module := range supportedModules {
				name := module.Name
				if module.Version != "" {
					name = fmt.Sprintf("%s:v%s", name, module.Version)
				}
				fmt.Println(name)
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getSupportedModulesCmd)
}
