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

// updateClusterCmd represents the cluster command
var updateClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "update cluster",
	Long:  `update cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterName, err := cmd.Flags().GetString("clustername")
		if err != nil {
			fmt.Println(err)
		}
		newClusterName, err := cmd.Flags().GetString("newclustername")
		if err != nil {
			fmt.Println(err)
		}
		newConfigFilePath, err := cmd.Flags().GetString("newconfigfilepath")
		if err != nil {
			fmt.Println(err)
		}

		if newClusterName == "" && newConfigFilePath == "" {
			fmt.Println("nothing to update")
		} else {
			_, err = api.PatchCluster(clusterName, newClusterName, newConfigFilePath)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("cluster updated successfully")
			}
		}
	},
}

func init() {
	updateClusterCmd.PersistentFlags().StringP("clustername", "", "", "")
	updateClusterCmd.PersistentFlags().StringP("newclustername", "", "", "")
	updateClusterCmd.PersistentFlags().StringP("newconfigfilepath", "", "", "")
	updateClusterCmd.MarkPersistentFlagRequired("clustername")
	updateCmd.AddCommand(updateClusterCmd)
}
