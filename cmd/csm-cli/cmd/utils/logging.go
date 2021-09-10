// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package utils

import (
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	// LogLevel - CLI log level environment variable
	LogLevel = "CSM_CLI_LOG_LEVEL"
)

var singletonLog *logrus.Logger
var once sync.Once

// GetLogger - get logger instance
func GetLogger() *logrus.Logger {
	once.Do(func() {
		singletonLog = logrus.New()
	})
	ChangeLogLevel(os.Getenv(LogLevel))
	return singletonLog
}

// ChangeLogLevel - updates the log level
func ChangeLogLevel(logLevel string) {

	switch strings.ToLower(logLevel) {

	case "debug":
		singletonLog.Level = logrus.DebugLevel
		break

	case "warn":
		singletonLog.Level = logrus.WarnLevel
		break

	case "error":
		singletonLog.Level = logrus.ErrorLevel
		break

	case "info":
		//Default level will be Info
		fallthrough

	default:
		singletonLog.Level = logrus.InfoLevel
	}
}
