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

// getSupportedStorageArraysCmd represents the supported-storage-arrays command
var getSupportedStorageArraysCmd = &cobra.Command{
	Use:   "supported-storage-arrays",
	Short: "get list of supported storage arrays",
	Long:  `get list of supported storage arrays`,
	Run: func(cmd *cobra.Command, args []string) {
		log := utils.GetLogger()
		supportedArrays, err := api.GetStorageTypes()
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("get supported storage arrays failed")
		} else if len(supportedArrays) == 0 {
			fmt.Println("no supported storage arrays")
		} else {
			for _, storageArray := range supportedArrays {
				fmt.Println(storageArray.Name)
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getSupportedStorageArraysCmd)
}
