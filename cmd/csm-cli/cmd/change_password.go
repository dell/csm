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

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api"
	utils "github.com/dell/csm-deployment/cmd/csm-cli/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// changePasswordCmd represents the password command
var changePasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "change password",
	Long:  `change password`,
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

		currentPassword, err := cmd.Flags().GetString("current-password")
		if err != nil {
			fmt.Println(err)
		}
		if currentPassword == "" {
			fmt.Println("Enter user's current password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Error("failed to read user's current password")
				os.Exit(0)
			} else {
				currentPassword = strings.TrimSpace(string(bytePassword))
			}
		}
		newPassword, err := cmd.Flags().GetString("new-password")
		if err != nil {
			fmt.Println(err)
		}
		if newPassword == "" {
			fmt.Println("Enter user's new password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Error("failed to read user's current password")
				os.Exit(0)
			} else {
				newPassword = strings.TrimSpace(string(bytePassword))
			}
		}

		err = api.ChangePassword(username, currentPassword, newPassword)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("change password failed")
		} else {
			fmt.Println("password has been changed")
		}
		err = api.LoginUser(username, newPassword)
		if err != nil {
			log.Errorf("error: %v", err)
			fmt.Println("authentication failed")
		} else {
			fmt.Println("authenticated successfully")
		}
	},
}

func init() {
	changePasswordCmd.PersistentFlags().StringP("username", "", "", "user's name")
	changePasswordCmd.PersistentFlags().StringP("current-password", "", "", "user's current password, skip to prompt for password")
	changePasswordCmd.PersistentFlags().StringP("new-password", "", "", "user's new password, skip to prompt for password")
	changePasswordCmd.MarkPersistentFlagRequired("username")
	changeCmd.AddCommand(changePasswordCmd)
}
