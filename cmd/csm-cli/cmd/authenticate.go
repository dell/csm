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

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "authenticate user",
	Long:  `authenticate user - login the user and set jwt`,
	Run: func(cmd *cobra.Command, args []string) {

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println(err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
		}

		_, err = api.LoginUser(username, password)
		if err != nil {
			fmt.Println(err)
			fmt.Println("authentication failed")
		} else {
			fmt.Println("authenticated successfully")
		}
	},
}

func init() {
	authenticateCmd.PersistentFlags().StringP("username", "u", "", "")
	authenticateCmd.PersistentFlags().StringP("password", "p", "", "")
	authenticateCmd.MarkPersistentFlagRequired("username")
	authenticateCmd.MarkPersistentFlagRequired("password")
	rootCmd.AddCommand(authenticateCmd)
}
