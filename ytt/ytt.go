package ytt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dell/csm-deployment/store"
	"github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
	"github.com/labstack/echo/v4"
)

type EchoLoggerWrapper struct {
	echoLogger echo.Logger
	debug      bool
}

func (e *EchoLoggerWrapper) Printf(s string, i ...interface{}) {
	e.echoLogger.Printf(s, i)
}

func (e *EchoLoggerWrapper) Debugf(s string, i ...interface{}) {
	if e.debug {
		e.echoLogger.Debugf(s, i)
	}
}

func (e *EchoLoggerWrapper) Warnf(s string, i ...interface{}) {
	e.echoLogger.Warnf(s, i)
}

func (e *EchoLoggerWrapper) DebugWriter() io.Writer {
	if e.debug {
		return os.Stderr
	}
	return noopWriter{}
}

type noopWriter struct{}

var _ io.Writer = noopWriter{}

func (w noopWriter) Write(data []byte) (int, error) { return len(data), nil }

type Output struct {
	*template.Output
}

func (out *Output) CreateAt(outputPath string) error {
	for _, file := range out.Files {
		err := file.Create(outputPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (out *Output) AsBytes() [][]byte {
	var res [][]byte
	for _, file := range out.Files {
		res = append(res, file.Bytes())
	}
	return res
}

func (out *Output) AsCombinedBytes() []byte {
	var res []byte
	for _, file := range out.Files {
		res = append(res, file.Bytes()...)
		res = append(res, []byte("---\n")...)
	}
	return res
}

type Interface interface {
	Template([]string, []string) (Output, error)
	TemplateFromApplication(appID uint,
		as store.ApplicationStateChangeStoreInterface, cs store.ClusterStoreInterface) (Output, error)
	SecretTemplateFromApplication(appID uint,
		as store.ApplicationStateChangeStoreInterface) (Output, error)
	NamespaceTemplateFromApplication(appID uint,
		as store.ApplicationStateChangeStoreInterface) (Output, error)
	GetEmptySecret(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error)
}

type client struct {
	logger       ui.UI
	templatePath string
}

type Option func(*client)

func WithLogger(logger echo.Logger, debug bool) Option {
	return func(c *client) {
		c.logger = &EchoLoggerWrapper{echoLogger: logger, debug: debug}
	}
}

func WithTemplatePath(path string) Option {
	return func(c *client) {
		c.templatePath = path
	}
}

var defaultOptions Option = func(c *client) {
	c.logger = ui.NewTTY(false)
	c.templatePath = "./"
}

func NewClient(opts ...Option) Interface {
	c := &client{}
	defaultOptions(c)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *client) Template(paths []string, values []string) (Output, error) {
	yttFiles, err := files.NewSortedFilesFromPaths(paths, files.SymlinkAllowOpts{})
	if err != nil {
		return Output{}, err
	}
	in := template.Input{Files: yttFiles}

	opts := template.NewOptions()
	opts.DataValuesFlags = template.DataValuesFlags{KVsFromStrings: values}

	output := opts.RunWithFiles(in, c.logger)

	return Output{&output}, output.Err
}

func (c *client) TemplateFromApplication(appID uint,
	as store.ApplicationStateChangeStoreInterface, cs store.ClusterStoreInterface) (Output, error) {
	c.logger.Printf("Generating template from app state with id %d \n", appID)
	appState, err := as.GetById(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	if len(appState.StorageArrays) == 0 {
		return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
	}
	arrayType := appState.StorageArrays[0].StorageArrayType
	c.logger.Printf("Array type is %q \n", arrayType.Name)

	cluster, err := cs.GetByID(appState.ClusterID)
	if err != nil {
		return Output{}, err
	}

	// only standalone modules can be installed without a driver
	if appState.DriverTypeID == 0 {
		for _, module := range appState.ModuleTypes {
			if !module.Standalone {
				return Output{}, fmt.Errorf("unable to install module %s without specifying a driver", module.Name)
			}
		}
	}

	paths := make([]string, 0)

	// if driver is going to be installed, load appropriate template files
	if appState.DriverTypeID > 0 {
		paths = append(paths, []string{
			c.templatePath + "templates/controller.yaml",
			c.templatePath + "templates/node.yaml",
			c.templatePath + "templates/csidriver.yaml",
			c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
			c.templatePath + "templates/common/values.yaml",
			c.templatePath + fmt.Sprintf("templates/common/k8s-%s-values.yaml", cluster.K8sVersion),
		}...)
	}

	// load standalone modules (ex: observability) and other shared ytt module files
	paths = append(paths, []string{
		c.templatePath + "templates/modules/",
		c.templatePath + "templates/configs/observability-0.3.0-values.yaml",
		c.templatePath + "templates/observability-0.3.0.yaml",
	}...)

	var values []string
	for _, module := range appState.ModuleTypes {
		values = append(values, fmt.Sprintf("%s.enabled=true", module.Name))
	}

	// add configuration values for the driver and modules that were passed from the API to create the application
	// these are space-delimited strings of the format "parent1.key1=value1 parent2.key2=value2"
	if len(appState.DriverConfiguration) > 0 {
		driverValues := strings.Split(appState.DriverConfiguration, " ")
		values = append(values, driverValues...)
	}

	if len(appState.ModuleConfiguration) > 0 {
		moduleValues := strings.Split(appState.ModuleConfiguration, " ")
		values = append(values, moduleValues...)
	}

	output, err := c.Template(paths, values)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func (c *client) SecretTemplateFromApplication(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error) {
	c.logger.Printf("Generating secret template from app state with id %d \n", appID)
	appState, err := as.GetById(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	if len(appState.StorageArrays) == 0 {
		return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
	}
	arrayType := appState.StorageArrays[0].StorageArrayType
	c.logger.Printf("Array type is %q \n", arrayType.Name)

	paths := []string{
		c.templatePath + "templates/driversecret.yaml",
		c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
		c.templatePath + fmt.Sprintf("templates/configs/values-%s-secret.yaml", arrayType.Name),
		c.templatePath + "templates/modules/driver_secret.lib.yml",
	}

	output, err := c.Template(paths, nil)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func (c *client) GetEmptySecret(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error) {
	c.logger.Printf("Generating secret template from app state with id %d \n", appID)
	appState, err := as.GetById(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	if len(appState.StorageArrays) == 0 {
		return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
	}
	arrayType := appState.StorageArrays[0].StorageArrayType
	c.logger.Printf("Array type is %q \n", arrayType.Name)
	paths := []string{
		c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
		c.templatePath + "templates/empty-secret.yaml",
	}

	output, err := c.Template(paths, nil)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func (c *client) NamespaceTemplateFromApplication(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error) {
	c.logger.Printf("Generating namespace template from app state with id %d \n", appID)
	appState, err := as.GetById(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	if len(appState.StorageArrays) == 0 {
		return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
	}
	arrayType := appState.StorageArrays[0].StorageArrayType

	paths := []string{
		c.templatePath + "templates/namespace.yaml",
		c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
	}

	output, err := c.Template(paths, nil)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}
