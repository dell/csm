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

// getTaskCmd represents the getTask command
var getTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "get tasks for applications",
	Long:  `get tasks for applications`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
		}

		applicationName, err := cmd.Flags().GetString("applicationname")
		if err != nil {
			fmt.Println(err)
		}

		if all {
			getTaskResp, err := api.GetAllTasks()
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get tasks failed")
			} else if len(getTaskResp) == 0 {
				fmt.Println("no tasks found")
			} else {
				for _, task := range getTaskResp {
					fmt.Println(fmt.Sprint(task.ID) + " " + task.Status + " " + task.ApplicationName)
				}
			}
		} else if applicationName != "" {
			getTaskResp, err := api.GetTaskByApplicationName(applicationName)
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get tasks failed")
			} else if len(getTaskResp) == 0 {
				fmt.Println("no task found for application: " + applicationName)
			} else {
				fmt.Println(fmt.Sprint(getTaskResp[0].ID) + " " + getTaskResp[0].Status + " " + getTaskResp[0].ApplicationName)
			}
		}
	},
}

func init() {
	getTaskCmd.PersistentFlags().StringP("applicationname", "", "", " get task for application with this name")
	getTaskCmd.Flags().BoolP("all", "A", false, "get all tasks")
	getCmd.AddCommand(getTaskCmd)
}
