// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package kapp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/dell/csm-deployment/utils"
)

// Interface is used to define the interface for kapp Interface
//go:generate mockgen -destination=mocks/kapp_interface.go -package=mocks github.com/dell/csm-deployment/kapp Interface
type Interface interface {
	DeployFromBytes(ctx context.Context, bytes []byte, appName string, wait bool, configPath string) (string, error)
	GetDeployDiff(ctx context.Context, bytes []byte, appName string, configPath string) (string, error)
	DeployFromFile(ctx context.Context, filePath string, appName string, configPath string) (string, error)
	Delete(ctx context.Context, appName string, configPath string) (string, error)
	List(ctx context.Context, namespace string, configPath string) (string, error)
}

type client struct {
	kappPath string
}

// NewClient -  - returns a new instance of client
func NewClient(kappPath string) Interface {
	if kappPath == "" {
		kappPath = utils.GetEnv("KAPP_BINARY", "/root/kapp")
	}

	return client{
		kappPath: kappPath,
	}
}

func (c client) DeployFromBytes(ctx context.Context, data []byte, appName string, wait bool, configPath string) (string, error) {
	r := bytes.NewReader(data)
	args := []string{"--kubeconfig", configPath, "deploy", "-a", appName, "-f", "-", "--yes", "--json"}
	if !wait {
		args = append(args, "--wait=false")
	}
	cmd := exec.CommandContext(ctx, c.kappPath, args...)
	cmd.Stdin = r
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", out), err
	}
	return fmt.Sprintf("%s", out), nil
}

func (c client) GetDeployDiff(ctx context.Context, data []byte, appName string, configPath string) (string, error) {
	r := bytes.NewReader(data)
	cmd := exec.CommandContext(ctx, c.kappPath, "--kubeconfig", configPath, "deploy", "-a", appName, "-f", "-", "--diff-run", "--json")
	cmd.Stdin = r
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", out), err
	}
	return fmt.Sprintf("%s", out), nil
}

func (c client) DeployFromFile(ctx context.Context, filePath string, appName string, configPath string) (string, error) {
	cmd := exec.CommandContext(ctx, c.kappPath, "--kubeconfig", configPath, "deploy", "-a", appName, "-f", filePath, "--yes", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", out), err
	}
	return fmt.Sprintf("%s", out), nil
}

func (c client) Delete(ctx context.Context, appName string, configPath string) (string, error) {
	cmd := exec.CommandContext(ctx, c.kappPath, "--kubeconfig", configPath, "delete", "-a", appName, "--yes", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", out), err
	}
	return fmt.Sprintf("%s", out), nil
}

func (c client) List(ctx context.Context, namespace string, configPath string) (string, error) {
	args := []string{"--kubeconfig", configPath, "list"}
	if namespace == "" {
		// list all namespaces
		args = append(args, "-A")
	} else {
		args = append(args, "-n", namespace)
	}
	cmd := exec.CommandContext(ctx, c.kappPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", out), err
	}
	return fmt.Sprintf("%s", out), nil
}
