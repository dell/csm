// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	api "github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// updateStorageCmd represents the storage command
var updateStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "update storage array",
	Long:  `update storage array`,
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
		newUniqueID, err := cmd.Flags().GetString("new-unique-id")
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
			fmt.Println("Do you want to update storage array's password:")
			reader := bufio.NewReader(os.Stdin)
			update, err := reader.ReadString('\n')
			if err != nil {
				log.Error("failed to read user's input")
				os.Exit(0)
			} else if strings.ToLower(update) == "y" || strings.ToLower(update) == "yes" {
				fmt.Println("Enter storage array's new password: ")
				bytePassword, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					log.Error("failed to read storage array's new password")
					os.Exit(0)
				} else {
					password = strings.TrimSpace(string(bytePassword))
				}
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

		_, err = api.PatchStorage(endpoint, username, password, uniqueID, storageType, newUniqueID, metaDataList)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("update storage array failed")
		} else {
			fmt.Println("storage array updated successfully")
		}
	},
}

func init() {
	updateStorageCmd.PersistentFlags().StringP("endpoint", "", "", "updated endpoint")
	updateStorageCmd.PersistentFlags().StringP("storage-type", "", "", "updated storage type")
	updateStorageCmd.PersistentFlags().StringP("unique-id", "", "", "unique id of storage to be updated")
	updateStorageCmd.PersistentFlags().StringP("new-unique-id", "", "", "updated unique id")
	updateStorageCmd.PersistentFlags().StringP("username", "", "", "updated username")
	updateStorageCmd.PersistentFlags().StringP("password", "", "", "updated password, skip to prompt for password")
	updateStorageCmd.PersistentFlags().StringP("meta-data", "", "", "comma separated list of meta-data")
	updateStorageCmd.MarkPersistentFlagRequired("unique-id")
	updateCmd.AddCommand(updateStorageCmd)
}
