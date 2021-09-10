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

// getApplicationCmd represents the application command
var getApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "get application",
	Long:  `get application`,
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

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}

		if all {
			getApplicationResp, err := api.GetAllApplications()
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get applications failed")
			} else if len(getApplicationResp) == 0 {
				fmt.Println("no applications created")
			} else {
				for _, application := range getApplicationResp {
					fmt.Println(application.Name + " " + application.Status)
				}
			}
		} else if name != "" {
			getApplicationResp, err := api.GetApplicationByName(name)
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get application failed")
			} else if len(getApplicationResp) == 0 {
				fmt.Println("unable to find application: " + name)
			} else {
				fmt.Println(getApplicationResp[0].Name + " " + getApplicationResp[0].Status)
			}
		}
	},
}

func init() {
	getApplicationCmd.PersistentFlags().StringP("name", "", "", "application's name")
	getApplicationCmd.Flags().BoolP("all", "A", false, "get all applications")
	getCmd.AddCommand(getApplicationCmd)
}
