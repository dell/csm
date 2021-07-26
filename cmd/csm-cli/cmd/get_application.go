// Package cmd
//Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getApplicationCmd represents the application command
var getApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "get application",
	Long:  `get application`,
	Run: func(cmd *cobra.Command, args []string) {
		//@TODO to be implemented
		fmt.Println("application get")
	},
}

func init() {
	getCmd.AddCommand(getApplicationCmd)
}
