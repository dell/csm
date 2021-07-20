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

	api "github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"

	"github.com/spf13/cobra"
)

// addClusterCmd represents the cluster command
var addClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "adds the cluster to CSM",
	Long:  `adds the cluster to CSM`,
	Run: func(cmd *cobra.Command, args []string) {
		clustername, err := cmd.Flags().GetString("clustername")
		if err != nil {
			fmt.Println(err)
		}

		configfilepath, err := cmd.Flags().GetString("configfilepath")
		if err != nil {
			fmt.Println(err)
		}

		_, err = api.AddCluster(clustername, configfilepath)
		if err != nil {
			fmt.Println(err)
			fmt.Println("add cluster failed")
		} else {
			fmt.Println("cluster added")
		}
	},
}

func init() {
	addClusterCmd.PersistentFlags().StringP("clustername", "", "", "")
	addClusterCmd.PersistentFlags().StringP("configfilepath", "", "", "")
	addClusterCmd.MarkPersistentFlagRequired("clustername")
	addClusterCmd.MarkPersistentFlagRequired("configfilepath")
	addCmd.AddCommand(addClusterCmd)
}
