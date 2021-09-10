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

// rejectTaskCmd represents the rejectTask command
var rejectTaskCmd = &cobra.Command{
	Use:   "reject-task",
	Short: "reject task for an application",
	Long:  `reject task for an application`,
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
		err = api.RejectTask(applicationName, update)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("reject task failed")
		} else {
			fmt.Println("task rejected successfully")
		}
	},
}

func init() {
	rejectTaskCmd.PersistentFlags().StringP("applicationname", "", "", "name of application corresponding to task")
	rejectTaskCmd.MarkPersistentFlagRequired("applicationname")
	rejectTaskCmd.Flags().BoolP("update", "u", false, "if task is to update application, then set this flag")
	rootCmd.AddCommand(rejectTaskCmd)
}
