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
	"strings"
	"syscall"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// addStorageCmd represents the storage command
var addStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "add storage array",
	Long:  `add storage array`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		endpoint, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			fmt.Println(err)
		}
		storageType, err := cmd.Flags().GetString("storage-type")
		if err != nil {
			fmt.Println(err)
		}
		uniqueID, err := cmd.Flags().GetString("unique-id")
		if err != nil {
			fmt.Println(err)
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println(err)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
		}
		if password == "" {
			fmt.Println("Enter storage array's password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Error("failed to read storage array's password")
				os.Exit(0)
			} else {
				password = strings.TrimSpace(string(bytePassword))
			}
		}

		metaData, err := cmd.Flags().GetString("meta-data")
		if err != nil {
			fmt.Println(err)
		}
		metaDataList := make([]string, 0)
		if metaData != "" {
			metaDataList = strings.Split(metaData, ",")
		}

		_, err = api.AddStorage(endpoint, username, password, uniqueID, storageType, metaDataList)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("add storage array failed")
		} else {
			fmt.Println("storage array added successfully")
		}
	},
}

func init() {
	addStorageCmd.PersistentFlags().StringP("endpoint", "", "", "storage array's endpoint")
	addStorageCmd.PersistentFlags().StringP("storage-type", "", "", "storage array type")
	addStorageCmd.PersistentFlags().StringP("unique-id", "", "", "storage array unique id")
	addStorageCmd.PersistentFlags().StringP("username", "", "", "storge array user's name")
	addStorageCmd.PersistentFlags().StringP("password", "", "", "storage array user's password, skip to prompt for password")
	addStorageCmd.PersistentFlags().StringP("meta-data", "", "", "comma separated list of meta-data")
	addStorageCmd.MarkPersistentFlagRequired("endpoint")
	addStorageCmd.MarkPersistentFlagRequired("storage-type")
	addStorageCmd.MarkPersistentFlagRequired("unique-id")
	addStorageCmd.MarkPersistentFlagRequired("username")
	addCmd.AddCommand(addStorageCmd)
}
