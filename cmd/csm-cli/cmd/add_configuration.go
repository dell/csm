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

// addConfigurationCmd represents the addconfiguration command
var addConfigurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "add configuration file",
	Long:  `add configuration file`,
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

		filePath, err := cmd.Flags().GetString("filepath")
		if err != nil {
			fmt.Println(err)
		}

		_, err = api.AddConfiguration(fileName, filePath)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("add configuration failed")
		} else {
			fmt.Println("configuration added")
		}
	},
}

func init() {
	addConfigurationCmd.PersistentFlags().StringP("filename", "", "", "configuration file name")
	addConfigurationCmd.PersistentFlags().StringP("filepath", "", "", "path to configuration file")
	addConfigurationCmd.MarkPersistentFlagRequired("filename")
	addConfigurationCmd.MarkPersistentFlagRequired("filepath")
	addCmd.AddCommand(addConfigurationCmd)
}
