// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SeverityEnum - Placeholder for Error Severity
type SeverityEnum string //@name SeverityEnum

const (
	// WarningSeverity captures enum value "WARNING"
	WarningSeverity SeverityEnum = "WARNING"
	// ErrorSeverity captures enum value "ERROR"
	ErrorSeverity SeverityEnum = "ERROR"
	// InfoSeverity captures enum value "INFO"
	InfoSeverity SeverityEnum = "INFO"
	// CriticalSeverity captures enum value "CRITICAL"
	CriticalSeverity SeverityEnum = "CRITICAL"
)

// HTTPStatusEnum Possible HTTP status values of completed or failed jobs.
type HTTPStatusEnum int32 //@name HTTPStatusEnum

// ErrorMessage - A message describing the failure, a contributing factor to the failure, or possibly the aftermath of the failure.
type ErrorMessage struct {
	// HTTPStatusEnum  Possible HTTP status values of completed or failed jobs
	ErrorCode HTTPStatusEnum `json:"code,omitempty" enums:"200,201,202,204,400,401,403,404,422,429,500,503"`
	// Message string.
	Message string `json:"message"`
	// Localized message
	MessageL10N interface{} `json:"message_l10n" swaggertype:"primitive,object"`
	Arguments   []string

	// SeverityEnum - The severity of the condition
	// * INFO - Information that may be of use in understanding the failure. It is not a problem to fix.
	// * WARNING - A condition that isn't a failure, but may be unexpected or a contributing factor. It may be necessary to fix the condition to successfully retry the request.
	// * ERROR - An actual failure condition through which the request could not continue.
	// * CRITICAL - A failure with significant impact to the system. Normally failed commands roll back and are just ERROR, but this is possible
	//
	Severity SeverityEnum `json:"severity,omitempty" enums:"INFO,WARNING,ERROR,CRITICAL"`
} //@name ErrorMessage

// ErrorResponse A standard response body used for all non-2xx REST responses.
type ErrorResponse struct {
	// HTTPStatusEnum  Possible HTTP status values of completed or failed jobs
	ErrorCode HTTPStatusEnum `json:"http_status_code,omitempty" enums:"200,201,202,204,400,401,403,404,422,429,500,503"`

	// A list of messages describing the failure encountered by this request. At least one will
	// be of Error severity because Info and Warning conditions do not cause the request to fail
	//
	Messages []*ErrorMessage `json:"messages"`
} //@name ErrorResponse

//BuildErrorMessage - Builds error message
func BuildErrorMessage(verbose string, code HTTPStatusEnum, severity SeverityEnum, errInterface interface{}) *ErrorMessage {
	errorResponse := ErrorResponse{}

	switch r := errInterface.(type) {
	case error:
		switch v := r.(type) {
		case *echo.HTTPError:
			cd := HTTPStatusEnum(v.Code)
			return &ErrorMessage{Message: verbose, MessageL10N: v.Message, ErrorCode: cd, Severity: SeverityEnum(severity)}
		default:
			return &ErrorMessage{Message: verbose, MessageL10N: v.Error(), ErrorCode: code, Severity: SeverityEnum(severity)}
		}

	case *http.Response:
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&errorResponse)
		cd := HTTPStatusEnum(r.StatusCode)
		if err != nil || errorResponse.Messages == nil {
			errMsg := verbose
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(dec.Buffered()); err == nil {
				s := buf.String()
				errMsg = fmt.Sprintf("%s: %s", errMsg, s)
			}
			return &ErrorMessage{ErrorCode: cd, Severity: severity,
				Message: errMsg}
		}
		firstErrMsg := (errorResponse.Messages)[0]
		firstErrMsg.ErrorCode = cd
		return firstErrMsg
	}
	return &ErrorMessage{Message: verbose, MessageL10N: errInterface, ErrorCode: code, Severity: SeverityEnum(severity)}
}

//BuildErrorResponse - Builds response with error
func BuildErrorResponse(code HTTPStatusEnum, severity SeverityEnum, verbose string, errInterface interface{}) ErrorResponse {
	e := ErrorResponse{}
	e.ErrorCode = code
	e.Messages = []*ErrorMessage{}
	e.Messages = append(e.Messages, BuildErrorMessage(verbose, code, severity, errInterface))
	return e
}

//NewErrorResponse - Returns response with error
func NewErrorResponse(code int, sve SeverityEnum, verbose string, err error) ErrorResponse {
	if verbose == "" {
		verbose = http.StatusText(code)
	}
	return BuildErrorResponse(HTTPStatusEnum(code), sve, verbose, err)
}
