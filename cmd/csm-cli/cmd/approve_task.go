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

// approveTaskCmd represents the approveTask command
var approveTaskCmd = &cobra.Command{
	Use:   "approve-task",
	Short: "approve task for application",
	Long:  `approve task for application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		applicationName, err := cmd.Flags().GetString("applicationname")
		if err != nil {
			fmt.Println(err)
		}
		update, err := cmd.Flags().GetBool("update")
		if err != nil {
			fmt.Println(err)
		}
		err = api.ApproveTask(applicationName, update)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("approve task failed")
		} else {
			fmt.Println("task approved successfully")
		}
	},
}

func init() {
	approveTaskCmd.PersistentFlags().StringP("applicationname", "", "", "name of application corresponding to task")
	approveTaskCmd.MarkPersistentFlagRequired("applicationname")
	approveTaskCmd.Flags().BoolP("update", "u", false, "if task is to update application, then set this flag")
	rootCmd.AddCommand(approveTaskCmd)
}
