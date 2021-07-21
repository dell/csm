package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

func LoginUser(username, password string) error {
	userLogin := &types.User{
		Username: username,
		Password: password,
	}

	err := saveAuthCreds(userLogin)
	if err != nil {
		return err
	}

	userLoginResponse := types.JWTToken
	err = HttpClient(http.MethodPost, UserLoginURI, nil, &userLoginResponse)
	if err != nil {
		return err
	}
	fmt.Println("Login success. Token: ", userLoginResponse)
	userLogin.Token = userLoginResponse
	err = saveAuthCreds(userLogin)
	if err != nil {
		return err
	}
	return nil
}

func saveAuthCreds(userLogin *types.User) error {
	file, _ := json.MarshalIndent(userLogin, "", " ")

	err := ioutil.WriteFile(filepath.Join(os.Getenv("AUTH_CONFIG_PATH"), "user.json"), file, 0600)
	if err != nil {
		return fmt.Errorf("failed to set user auth creds with error %v", err)
	}
	return nil
}
