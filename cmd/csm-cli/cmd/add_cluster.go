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
	"log"
	"os"

	api "github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"

	"github.com/spf13/cobra"
)

// addClusterCmd represents the cluster command
var addClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "adds the cluster to CSM",
	Long:  `adds the cluster to CSM`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
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
			log.Errorf("error: %v", err)
			fmt.Println("add cluster failed")
		} else {
			fmt.Println("cluster added")
		}
	},
}

func init() {
	addClusterCmd.PersistentFlags().StringP("clustername", "", "", "cluster's name")
	addClusterCmd.PersistentFlags().StringP("configfilepath", "", "", "path to cluster's kube config file")
	err := addClusterCmd.MarkPersistentFlagRequired("clustername")
	if err != nil {
		log.Fatal(err)
	}
	err = addClusterCmd.MarkPersistentFlagRequired("configfilepath")
	if err != nil {
		log.Fatal(err)
	}
	addCmd.AddCommand(addClusterCmd)
}
