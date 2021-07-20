package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var ApiServer = fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME"), os.Getenv("HOST"), os.Getenv("PORT"))

func GetClient(protocol string) *http.Client {
	if protocol == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	} else {
		return &http.Client{}
	}
}

func DoAndGetResponse(httpReq *http.Request, client *http.Client, resp interface{}) error {
	res, err := client.Do(httpReq)
	if err != nil {
		return errors.New(fmt.Sprintf("request failed with error: %v", err))
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("request failed with status code: %d", res.StatusCode))
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(resp)
	if err != nil {
		return errors.New(fmt.Sprintf("decoding the response failed with error: %v", err))
	}
	return nil
}

func GetJWTToken() (string, error) {
	token, err := ioutil.ReadFile(filepath.Join(os.Getenv("JWTPATH"), "jwt"))
	if err != nil {
		return "", errors.New("unable to retrieve jwt token")
	}
	return string(token), nil
}

func HttpClient(method, uri string, req, resp interface{}) error {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to parse the request JSON with error: %v", err))
	}

	protocol := os.Getenv("SCHEME")
	url := fmt.Sprintf("%s%s", ApiServer, uri)
	client := GetClient(protocol)

	jwtToken, err := GetJWTToken()
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(reqJson))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create request with error: %v", err))
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	httpReq.Header.Set("Content-Type", "application/json")

	err = DoAndGetResponse(httpReq, client, resp)
	if err != nil {
		return err
	}
	return nil
}

func HttpClusterClient(method, uri, configFilePath string, reqFields map[string]string, resp interface{}) error {
	body, writer, err := createClusterMultipartFormData(reqFields, configFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("creating formdata failed with error: %v", err))
	}

	protocol := os.Getenv("SCHEME")
	url := fmt.Sprintf("%s%s", ApiServer, uri)
	client := GetClient(protocol)

	jwtToken, err := GetJWTToken()
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(method, url, &body)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create request with error: %v", err))
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	err = DoAndGetResponse(httpReq, client, resp)
	if err != nil {
		return err
	}
	return nil
}

func createClusterMultipartFormData(reqFields map[string]string, configFilePath string) (bytes.Buffer, *multipart.Writer, error) {
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
