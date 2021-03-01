package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func ensureTargetConnection() {
	// If the accessToken is not set or testAccesToken returns false
	if accessToken == "" || testAccessToken() == false {
		var authError error
		// Authenticate
		if apiToken != "" {
			accessToken, authError = authenticateCloud(viper.GetString("target."+currentTargetName+".server"), viper.GetString("target."+currentTargetName+".apiToken"))
		} else {
			accessToken, authError = authenticateOnPrem(viper.GetString("target."+currentTargetName+".server"), viper.GetString("target."+currentTargetName+".username"), viper.GetString("target."+currentTargetName+".password"), viper.GetString("target."+currentTargetName+".domain"))
		}
		if authError != nil {
			fmt.Println("Authentication failed", authError.Error())
			os.Exit(1)
		}
		viper.Set("target."+currentTargetName+".accessToken", accessToken)
		viper.WriteConfig()
		accessToken = viper.GetString("target." + currentTargetName + ".accessToken")
	}
}

func authenticateOnPrem(server string, username string, password string, domain string) (string, error) {
	client := resty.New()
	response, err := client.R().
		SetBody(AuthenticationRequest{username, password, domain}).
		SetResult(&AuthenticationResponse{}).
		SetError(&AuthenticationError{}).
		Post("https://" + server + "/csp/gateway/am/idp/auth/login?access_token")
	if response.IsError() {
		return "", err
	}
	return response.Result().(*AuthenticationResponse).AccessToken, err
}
func authenticateCloud(server string, apiToken string) (string, error) {
	client := resty.New()
	response, err := client.R().
		SetBody(AuthenticationRequestCloud{apiToken}).
		SetResult(&AuthenticationResponseCloud{}).
		SetError(&AuthenticationError{}).
		Post("https://" + server + "/iaas/api/login")
	if response.IsError() {
		return "", err
	}
	return response.Result().(*AuthenticationResponseCloud).Token, err
}

func testAccessToken() bool {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(accessToken).
		Get("https://" + server + "/iaas/api/projects")
	if err != nil {
		return false
	}
	if response.StatusCode() == 401 {
		//fmt.Println("Token authentication failed: ", response.StatusCode())
		return false
	}
	return true
}
