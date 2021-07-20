package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

func LoginUser(username, password string) (*types.UserLoginResponse, error) {
	userLoginReq := &types.User{
		Username: username,
		Password: password,
	}

	userLoginResponse := &types.UserLoginResponse{}
	err := HttpClient(http.MethodPost, UserLoginURI, userLoginReq, userLoginResponse)
	if err != nil {
		return nil, err
	}
	fmt.Println("Login success. Token: ", userLoginResponse.Token)
	err = ioutil.WriteFile(filepath.Join(os.Getenv("JWTPATH"), "jwt"), []byte(userLoginResponse.Token), 0755)
	if err != nil {
		return nil, errors.New("unable to set jwt token")
	}
	return userLoginResponse, nil
}
