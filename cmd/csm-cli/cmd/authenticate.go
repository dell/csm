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
	"strings"
	"syscall"

	api "github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "authenticate user",
	Long:  `authenticate user - login the user and set jwt`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		log := utils.GetLogger()
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println(err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(err)
		}
		if password == "" {
			fmt.Println("Enter user's password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println("failed to read user's password")
			} else {
				password = strings.TrimSpace(string(bytePassword))
			}
		}

		err = api.LoginUser(username, password)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("authentication failed")
		} else {
			fmt.Println("authenticated successfully")
		}
	},
}

func init() {
	authenticateCmd.PersistentFlags().StringP("username", "u", "", "user's name")
	authenticateCmd.PersistentFlags().StringP("password", "p", "", "user's password, skip to prompt for password")
	authenticateCmd.MarkPersistentFlagRequired("username")
	rootCmd.AddCommand(authenticateCmd)
}
