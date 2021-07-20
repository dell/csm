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

// getClusterCmd represents the cluster command
var getClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "get clusters",
	Long:  `get clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		clustername, err := cmd.Flags().GetString("clustername")
		if err != nil {
			fmt.Println(err)
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
		}

		if clustername != "" {
			clusters, err := api.GetClusterByName(clustername)
			if err != nil {
				fmt.Println(err)
				fmt.Println("get cluster failed")
			}
			for _, cluster := range clusters {
				fmt.Println(fmt.Sprint(cluster.ClusterId) + " " + cluster.ClusterName + " " + cluster.Nodes)
			}
		} else if all {
			clusters, err := api.GetAllClusters()
			if err != nil {
				fmt.Println(err)
				fmt.Println("get cluster failed")
			} else if len(clusters) == 0 {
				fmt.Println("no clusters added")
			} else {
				for _, cluster := range clusters {
					fmt.Println(fmt.Sprint(cluster.ClusterId) + " " + cluster.ClusterName + " " + cluster.Nodes)
				}
			}
		} else {
			fmt.Println("check --help option for this command")
		}

	},
}

func init() {
	getClusterCmd.PersistentFlags().StringP("clustername", "", "", "")
	getClusterCmd.Flags().BoolP("all", "A", false, "get all clusters")
	getCmd.AddCommand(getClusterCmd)
}
