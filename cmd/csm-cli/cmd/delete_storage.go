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

// deleteStorageCmd represents the storage command
var deleteStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "delete storage array",
	Long:  `delete storage array`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		uniqueID, err := cmd.Flags().GetString("unique-id")
		if err != nil {
			fmt.Println(err)
		}
		err = api.DeleteStorage(uniqueID)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("delete storage array failed")
		} else {
			fmt.Println("storage array deleted successfully")
		}
	},
}

func init() {
	deleteStorageCmd.PersistentFlags().StringP("unique-id", "", "", "storage array's unique ID")
	deleteStorageCmd.MarkPersistentFlagRequired("unique-id")
	deleteCmd.AddCommand(deleteStorageCmd)
}
