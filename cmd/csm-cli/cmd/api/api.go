// Package api for API services
// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

// APIServer - Placeholder for API Server
var APIServer = fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME"), os.Getenv("API_SERVER_IP"), os.Getenv("API_SERVER_PORT"))

// GetClient - Get http/https client
func GetClient(protocol string) *http.Client {
	if protocol == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	}
	return &http.Client{}
}

// DoAndGetResponse - Execute http request and return response
func DoAndGetResponse(httpReq *http.Request, client *http.Client, resp interface{}) error {
	res, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed with error: %v", err)
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusNoContent {
		errorResp := types.ErrorResponse{}
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&errorResp)
		if err != nil && err != io.EOF && httpReq.Method != http.MethodDelete {
			return fmt.Errorf("decoding the response failed with error: %v", err)
		}
		if len(errorResp.Messages) > 0 {
			return fmt.Errorf("request failed with status code: %d. message: %s", res.StatusCode, errorResp.Messages[0].Message)
		}
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(resp)
	if err != nil && err != io.EOF && httpReq.Method != http.MethodDelete {
		return fmt.Errorf("decoding the response failed with error: %v", err)
	}
	return nil
}

// GetUserAuthCreds - Get user authentication from AUTH_CONFIG_PATH
func GetUserAuthCreds() (string, string, string, error) {
	authCredsFile, err := ioutil.ReadFile(filepath.Join(os.Getenv("AUTH_CONFIG_PATH"), "user.json"))
	if err != nil {
		return "", "", "", fmt.Errorf("unable to read user creds")
	}
	creds := types.User{}
	err = json.Unmarshal([]byte(authCredsFile), &creds)
	if err != nil {
		return "", "", "", fmt.Errorf("unable to parse user creds")
	}
	return Username, Password, creds.Token, nil
}

// basicAuth - encode basic auth creds
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// HTTPClient - Execute http request based on method and uri
func HTTPClient(method, uri string, req, resp interface{}) error {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to parse the request JSON with error: %v", err)
	}

	protocol := os.Getenv("SCHEME")
	url := fmt.Sprintf("%s%s", APIServer, uri)
	client := GetClient(protocol)

	username, password, jwtToken, err := GetUserAuthCreds()
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return fmt.Errorf("failed to create request with error: %v", err)
	}
	if username != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(username, password)))
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if jwtToken != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	}

	err = DoAndGetResponse(httpReq, client, resp)
	if err != nil {
		return err
	}
	return nil
}

// HTTPFormDataClient - Execute http request based on method uri and config file
func HTTPFormDataClient(method, uri, configFilePath string, reqFields map[string]string, resp interface{}) error {
	body, writer, err := createMultipartFormData(reqFields, configFilePath)
	if err != nil {
		return fmt.Errorf("creating formdata failed with error: %v", err)
	}

	protocol := os.Getenv("SCHEME")
	url := fmt.Sprintf("%s%s", APIServer, uri)
	client := GetClient(protocol)

	username, password, jwtToken, err := GetUserAuthCreds()
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(method, url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request with error: %v", err)
	}
	if username != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(username, password)))
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	err = DoAndGetResponse(httpReq, client, resp)
	if err != nil {
		return err
	}
	return nil
}

// createMultipartFormData - create form data request
func createMultipartFormData(reqFields map[string]string, configFilePath string) (bytes.Buffer, *multipart.Writer, error) {
	var body bytes.Buffer
	var err error
	writer := multipart.NewWriter(&body)
	if configFilePath != "" {
		var fw io.Writer
		file, err := os.Open(configFilePath)
		if err != nil {
			return body, nil, err
		}
		if fw, err = writer.CreateFormFile("file", file.Name()); err != nil {
			return body, nil, err
		}
		if _, err = io.Copy(fw, file); err != nil {
			return body, nil, err
		}
	}
	for key, value := range reqFields {
		if err = writer.WriteField(key, value); err != nil {
			return body, nil, err
		}
	}
	writer.Close()
	return body, writer, nil
}
