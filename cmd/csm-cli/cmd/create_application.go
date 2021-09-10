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
	"time"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	// DefaultAutoApproveRetries - Default retries for auto approve application
	DefaultAutoApproveRetries = 5
	// DefaultAutoApproveTimeout - Default timeout for aut approve application
	DefaultAutoApproveTimeout = 10 * time.Second
)

// createApplicationCmd represents the application command
var createApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "create application",
	Long:  `create application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		clustername, err := cmd.Flags().GetString("clustername")
		if err != nil {
			fmt.Println(err)
		}

		driverConfiguration, err := cmd.Flags().GetString("driver-configuration")
		if err != nil {
			fmt.Println(err)
		}
		driverConfigurationList := make([]string, 0)
		if driverConfiguration != "" {
			driverConfigurationList = strings.Split(driverConfiguration, ",")
		}

		driverType, err := cmd.Flags().GetString("driver-type")
		if err != nil {
			fmt.Println(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}

		storageArrays, err := cmd.Flags().GetString("storage-arrays")
		if err != nil {
			fmt.Println(err)
		}
		storageArraysList := strings.Split(storageArrays, ",")

		moduleConfiguration, err := cmd.Flags().GetString("module-configuration")
		if err != nil {
			fmt.Println(err)
		}
		moduleConfigurationList := make([]string, 0)
		if moduleConfiguration != "" {
			moduleConfigurationList = strings.Split(moduleConfiguration, ",")
		}

		moduleType, err := cmd.Flags().GetString("module-type")
		if err != nil {
			fmt.Println(err)
		}
		moduleTypeList := make([]string, 0)
		if moduleType != "" {
			moduleTypeList = strings.Split(moduleType, ",")
		}

		err = api.CreateApplication(name, clustername, driverType, driverConfigurationList, storageArraysList, moduleTypeList, moduleConfigurationList)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("create application failed")
		} else {
			fmt.Println("application created successfully")
		}

		autoApprove, err := cmd.Flags().GetBool("autoapprove")
		if err != nil {
			fmt.Println(err)
		}

		if autoApprove {
			for i := 0; i < DefaultAutoApproveRetries; i++ {
				log.Debugf("Retry %d for auto-approve", i)
				time.Sleep(DefaultAutoApproveTimeout)
				err := api.ApproveTask(name, false)
				if err != nil {
					log.Errorf("error: %v", err)
					fmt.Println("approving the task failed")
				} else {
					fmt.Println("task approved successfully")
					break
				}
			}
		}
	},
}

func init() {
	createApplicationCmd.PersistentFlags().StringP("clustername", "", "", "cluster name")
	createApplicationCmd.PersistentFlags().StringP("driver-configuration", "", "", "comma separated driver configuration parameters")
	createApplicationCmd.PersistentFlags().StringP("driver-type", "", "", "application's driver type")
	createApplicationCmd.PersistentFlags().StringP("name", "", "", "application name")
	createApplicationCmd.PersistentFlags().StringP("storage-arrays", "", "", "comma separated storage arrays unique id's")
	createApplicationCmd.PersistentFlags().StringP("module-type", "", "", "comma separated list of module types")
	createApplicationCmd.PersistentFlags().StringP("module-configuration", "", "", "comma separated module configuration parameters")
	createApplicationCmd.Flags().BoolP("autoapprove", "", false, "if set, approves the task created for application")
	createApplicationCmd.MarkPersistentFlagRequired("clustername")
	createApplicationCmd.MarkPersistentFlagRequired("driver-type")
	createApplicationCmd.MarkPersistentFlagRequired("name")
	createApplicationCmd.MarkPersistentFlagRequired("storage-arrays")
	createCmd.AddCommand(createApplicationCmd)
}
