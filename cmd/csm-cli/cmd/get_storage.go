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
	"github.com/spf13/cobra"
)

// getStorageCmd represents the storage command
var getStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "get storage arrays",
	Long:  `get storage arrays`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqueId, err := cmd.Flags().GetString("unique-id")
		if err != nil {
			fmt.Println(err)
		}
		if uniqueId != "" {
			storageArrays, err := api.GetStorageByParam(api.StorageTypeIdResponseField, uniqueId)
			if err != nil {
				fmt.Println(err)
				fmt.Println("get storage array by unique id failed")
			} else {
				fmt.Println("get storage array result for unique id " + uniqueId + ":")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueId)
				}
			}
		}
		endpoint, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			fmt.Println(err)
		}
		if endpoint != "" {
			storageArrays, err := api.GetStorageByParam(api.EndpointResponseField, endpoint)
			if err != nil {
				fmt.Println(err)
				fmt.Println("get storage array by endpoint failed")
			} else {
				fmt.Println("get storage array result for endpoint " + endpoint + ":")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueId)
				}
			}
		}
		storageType, err := cmd.Flags().GetString("storage-type")
		if err != nil {
			fmt.Println(err)
		}
		storageTypeId := api.GetStorageTypeId(storageType)
		if storageType != "" {
			storageArrays, err := api.GetStorageByParam(api.StorageTypeIdResponseField, storageTypeId)
			if err != nil {
				fmt.Println(err)
				fmt.Println("get storage array by storage type failed")
			} else {
				fmt.Println("get storage array result for storage type " + storageType + ":")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueId)
				}
			}
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
		}
		if all {
			storageArrays, err := api.GetAllStorage()
			if err != nil {
				fmt.Println(err)
				fmt.Println("get all storage arrays failed")
			} else {
				fmt.Println("get all storage arrays result:")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueId)
				}
			}
		}
	},
}

func init() {
	getStorageCmd.PersistentFlags().StringP("storage-type", "", "", "")
	getStorageCmd.PersistentFlags().StringP("unique-id", "", "", "")
	getStorageCmd.PersistentFlags().StringP("endpoint", "", "", "")
	getStorageCmd.Flags().BoolP("all", "A", false, "get all storage arrays")
	getCmd.AddCommand(getStorageCmd)
}
