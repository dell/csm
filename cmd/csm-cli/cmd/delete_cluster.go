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

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the cluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "delete a cluster",
	Long:  `delete a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		clustername, err := cmd.Flags().GetString("clustername")
		if err != nil {
			fmt.Println(err)
		}

		err = api.DeleteCluster(clustername)
		if err != nil {
			fmt.Println(err)
			fmt.Println("delete cluster failed")
		} else {
			fmt.Println("cluster deleted successfully")
		}
	},
}

func init() {
	deleteClusterCmd.PersistentFlags().StringP("clustername", "", "", "")
	deleteClusterCmd.MarkPersistentFlagRequired("clustername")
	deleteCmd.AddCommand(deleteClusterCmd)
}
