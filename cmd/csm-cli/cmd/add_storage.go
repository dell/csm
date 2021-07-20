/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	"github.com/spf13/cobra"
)

// addStorageCmd represents the storage command
var addStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "add storage array",
	Long:  `add storage array`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			fmt.Println(err)
		}
		storageType, err := cmd.Flags().GetString("storage-type")
		if err != nil {
			fmt.Println(err)
		}
		uniqueId, err := cmd.Flags().GetString("unique-id")
		if err != nil {
			fmt.Println(err)
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println(err)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
		}

		_, err = api.AddStorage(endpoint, username, password, uniqueId, storageType)
		if err != nil {
			fmt.Println(err)
			fmt.Println("add storage array failed")
		} else {
			fmt.Println("storage array added successfully")
		}
	},
}

func init() {
	addStorageCmd.PersistentFlags().StringP("endpoint", "", "", "")
	addStorageCmd.PersistentFlags().StringP("storage-type", "", "", "")
	addStorageCmd.PersistentFlags().StringP("unique-id", "", "", "")
	addStorageCmd.PersistentFlags().StringP("username", "", "", "")
	addStorageCmd.PersistentFlags().StringP("password", "", "", "")
	addStorageCmd.MarkPersistentFlagRequired("endpoint")
	addStorageCmd.MarkPersistentFlagRequired("storage-type")
	addStorageCmd.MarkPersistentFlagRequired("unique-id")
	addStorageCmd.MarkPersistentFlagRequired("username")
	addStorageCmd.MarkPersistentFlagRequired("password")
	addCmd.AddCommand(addStorageCmd)
}
