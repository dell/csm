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

	api "github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
)

// getSupportedDriversCmd represents the supported-drivers command
var getSupportedDriversCmd = &cobra.Command{
	Use:   "supported-drivers",
	Short: "get list of supported drivers",
	Long:  `get list of supported drivers`,
	Run: func(cmd *cobra.Command, args []string) {
		log := utils.GetLogger()
		supportedDrivers, err := api.GetDriverTypes()
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("get supported drivers failed")
		} else if len(supportedDrivers) == 0 {
			fmt.Println("no supported drivers")
		} else {
			for _, driver := range supportedDrivers {
				name := driver.StorageType
				if driver.Version != "" {
					name = fmt.Sprintf("%s:v%s", name, driver.Version)
				}
				fmt.Println(name)
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getSupportedDriversCmd)
}
