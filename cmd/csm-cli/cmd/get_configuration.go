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

// getConfigurationCmd represents the get configuration command
var getConfigurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "get configuration files",
	Long:  `get configuration files`,
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
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
		}

		if fileName != "" {
			configs, err := api.GetConfigurationByName(fileName)
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get configuration failed")
			}
			for _, config := range configs {
				fmt.Println(fmt.Sprint(config.ID) + " " + config.Name)
			}
		} else if all {
			configs, err := api.GetAllConfigurations()
			if err != nil {
				log.Errorf("error: %v", err)
				fmt.Println("get configuration failed")
			} else if len(configs) == 0 {
				fmt.Println("no configuration added")
			} else {
				for _, config := range configs {
					fmt.Println(fmt.Sprint(config.ID) + " " + config.Name)
				}
			}
		}
	},
}

func init() {
	getConfigurationCmd.PersistentFlags().StringP("filename", "", "", "configuration file name")
	getConfigurationCmd.Flags().BoolP("all", "A", false, "get all configuration files")
	getCmd.AddCommand(getConfigurationCmd)
}
