package kapp

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/dell/csm-deployment/utils"
)

type Interface interface {
	DeployFromBytes(ctx context.Context, bytes []byte, appName string, wait bool) (string, error)
	GetDeployDiff(ctx context.Context, bytes []byte, appName string) (string, error)
	DeployFromFile(ctx context.Context, filePath string, appName string) (string, error)
	Delete(ctx context.Context, appName string) (string, error)
	List(ctx context.Context, namespace string) (string, error)
}

type client struct {
	config   string
	kappPath string
}

func NewClient(kappPath string, configPath string) Interface {
	if kappPath == "" {
		kappPath = utils.GetEnv("KAPP_BINARY", "/root/kapp")
	}

	return client{
		kappPath: kappPath,
		config:   configPath,
	}
}

func (c client) DeployFromBytes(ctx context.Context, data []byte, appName string, wait bool) (string, error) {
	r := bytes.NewReader(data)
	args := []string{"--kubeconfig", c.config, "deploy", "-a", appName, "-f", "-", "--yes", "--json"}
	if !wait {
		args = append(args, "--wait=false")
	}
	cmd := exec.CommandContext(ctx, c.kappPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = r
	err := cmd.Run()
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func (c client) GetDeployDiff(ctx context.Context, data []byte, appName string) (string, error) {
	r := bytes.NewReader(data)
	cmd := exec.CommandContext(ctx, c.kappPath, "--kubeconfig", c.config, "deploy", "-a", appName, "-f", "-", "--diff-run", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = r
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (c client) DeployFromFile(ctx context.Context, filePath string, appName string) (string, error) {
	cmd := exec.CommandContext(ctx, c.kappPath, "--kubeconfig", c.config, "deploy", "-a", appName, "-f", filePath, "--yes", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (c client) Delete(ctx context.Context, appName string) (string, error) {
	cmd := exec.CommandContext(ctx, c.kappPath, "--kubeconfig", c.config, "delete", "-a", appName, "--yes", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (c client) List(ctx context.Context, namespace string) (string, error) {
	args := []string{"--kubeconfig", c.config, "list"}
	if namespace == "" {
		// list all namespaces
		args = append(args, "-A")
	} else {
		args = append(args, "-n", namespace)
	}
	cmd := exec.CommandContext(ctx, c.kappPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
