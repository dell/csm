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

	api "github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
)

// updateconfigurationCmd represents the updateconfiguration command
var updateConfigurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "update configuration file",
	Long:  `update configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		fileName, err := cmd.Flags().GetString("filename")
		if err != nil {
			fmt.Println(err)
		}
		newFileName, err := cmd.Flags().GetString("newfilename")
		if err != nil {
			fmt.Println(err)
		}
		newFilePath, err := cmd.Flags().GetString("newfilepath")
		if err != nil {
			fmt.Println(err)
		}

		if newFileName == "" && newFilePath == "" {
			fmt.Println("nothing to update")
		} else {
			_, err = api.PatchConfiguration(fileName, newFileName, newFilePath)
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("update configuration failed")
			} else {
				fmt.Println("configuration updated successfully")
			}
		}
	},
}

func init() {
	updateConfigurationCmd.PersistentFlags().StringP("filename", "", "", "configuration file name")
	updateConfigurationCmd.PersistentFlags().StringP("newfilename", "", "", "updated configuration file name")
	updateConfigurationCmd.PersistentFlags().StringP("newfilepath", "", "", "updated configuration file path")
	updateConfigurationCmd.MarkPersistentFlagRequired("filename")
	updateConfigurationCmd.MarkPersistentFlagRequired("newfilepath")
	updateCmd.AddCommand(updateConfigurationCmd)
}
