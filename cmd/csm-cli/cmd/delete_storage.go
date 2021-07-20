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

// deleteStorageCmd represents the storage command
var deleteStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "delete storage array",
	Long:  `delete storage array`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqueId, err := cmd.Flags().GetString("unique-id")
		if err != nil {
			fmt.Println(err)
		}
		err = api.DeleteStorage(uniqueId)
		if err != nil {
			fmt.Println(err)
			fmt.Println("delete storage array failed")
		} else {
			fmt.Println("storage array deleted successfully")
		}
	},
}

func init() {
	deleteStorageCmd.PersistentFlags().StringP("unique-id", "", "", "")
	deleteStorageCmd.MarkPersistentFlagRequired("unique-id")
	deleteCmd.AddCommand(deleteStorageCmd)
}
