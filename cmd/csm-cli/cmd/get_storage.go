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
	"os"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
)

// getStorageCmd represents the storage command
var getStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "get storage arrays",
	Long:  `get storage arrays`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		uniqueID, err := cmd.Flags().GetString("unique-id")
		if err != nil {
			fmt.Println(err)
		}
		if uniqueID != "" {
			storageArrays, err := api.GetStorageByParam(api.StorageUniqueIDResponseField, uniqueID)
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get storage array by unique id failed")
			} else if len(storageArrays) == 0 {
				fmt.Println("no storage array with unique id " + uniqueID)
			} else {
				fmt.Println("get storage array result for unique id " + uniqueID + ":")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueID + " " + array.Username)
				}
			}
		}
		storageType, err := cmd.Flags().GetString("storage-type")
		if err != nil {
			fmt.Println(err)
		}
		if storageType != "" {
			storageArrays, err := api.GetStorageByParam(api.StorageTypeResponseField, storageType)
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get storage array by storage type failed")
			} else if len(storageArrays) == 0 {
				fmt.Println("no storage array with storage type " + storageType)
			} else {
				fmt.Println("get storage array result for storage type " + storageType + ":")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueID + " " + array.Username)
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
			} else if len(storageArrays) == 0 {
				fmt.Println("no storage arrays added")
			} else {
				fmt.Println("get all storage arrays result:")
				for _, array := range storageArrays {
					fmt.Println(array.Endpoint + " " + array.UniqueID + " " + array.Username)
				}
			}
		}
	},
}

func init() {
	getStorageCmd.PersistentFlags().StringP("storage-type", "", "", "storage array's type")
	getStorageCmd.PersistentFlags().StringP("unique-id", "", "", "storage array's unique ID")
	getStorageCmd.Flags().BoolP("all", "A", false, "get all storage arrays")
	getCmd.AddCommand(getStorageCmd)
}
