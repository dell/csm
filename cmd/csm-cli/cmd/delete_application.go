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

// deleteApplicationCmd represents the application command
var deleteApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "delete application",
	Long:  `delete application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}

		err = api.DeleteApplication(name)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("delete application failed")
		} else {
			fmt.Println("application deleted successfully")
		}
	},
}

func init() {
	deleteApplicationCmd.PersistentFlags().StringP("name", "", "", "application name")
	deleteApplicationCmd.MarkPersistentFlagRequired("name")
	deleteCmd.AddCommand(deleteApplicationCmd)
}
