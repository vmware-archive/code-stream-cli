package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func ensureTargetConnection() {
	// If the targetConfig.accesstoken is not set or testAccesToken returns false
	if targetConfig.accesstoken == "" || testAccessToken() == false {
		var authError error
		// Authenticate
		if targetConfig.apitoken != "" {
			targetConfig.accesstoken, authError = authenticateCloud(targetConfig)
		} else {
			targetConfig.accesstoken, authError = authenticateOnPrem(targetConfig)
		}
		if authError != nil {
			fmt.Println("Authentication failed", authError.Error())
			os.Exit(1)
		}
		viper.Set("target."+currentTargetName+".targetConfig.accesstoken", targetConfig.accesstoken)
		viper.WriteConfig()
		targetConfig.accesstoken = viper.GetString("target." + currentTargetName + ".targetConfig.accesstoken")
	}
}

func authenticateOnPrem(target config) (string, error) {
	client := resty.New()
	response, err := client.R().
		SetBody(AuthenticationRequest{target.username, target.password, target.domain}).
		SetResult(&AuthenticationResponse{}).
		SetError(&AuthenticationError{}).
		Post("https://" + target.server + "/csp/gateway/am/idp/auth/login?access_token")
	if response.IsError() {
		return "", err
	}
	return response.Result().(*AuthenticationResponse).AccessToken, err
}
func authenticateCloud(target config) (string, error) {
	client := resty.New()
	response, err := client.R().
		SetBody(AuthenticationRequestCloud{target.apitoken}).
		SetResult(&AuthenticationResponseCloud{}).
		SetError(&AuthenticationError{}).
		Post("https://" + target.server + "/iaas/api/login")
	if response.IsError() {
		return "", err
	}
	return response.Result().(*AuthenticationResponseCloud).Token, err
}

func testAccessToken() bool {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/iaas/api/projects")
	if err != nil {
		return false
	}
	if response.StatusCode() == 401 {
		//fmt.Println("Token authentication failed: ", response.StatusCode())
		return false
	}
	return true
}

func exportYaml(name, project, path, object string) {
	var exportPath string
	var qParams = make(map[string]string)
	qParams[object] = name
	qParams["project"] = project
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		SetAuthToken(targetConfig.accesstoken).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		Get("https://" + targetConfig.server + "/pipeline/api/export")

	if queryResponse.IsError() {
		fmt.Println("Export failed", err)
		os.Exit(1)
	}
}

// importYaml import a yaml pipeline or endpoint
func importYaml(yamlPath, action string) bool {
	var qParams = make(map[string]string)
	qParams["action"] = action
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}
	yamlPayload := string(yamlBytes)
	client := resty.New()
	response, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/x-yaml").
		SetBody(yamlPayload).
		SetAuthToken(targetConfig.accesstoken).
		Post("https://" + targetConfig.server + "/pipeline/api/import")
	if response.IsError() {
		fmt.Println("Import/Update failed", response.StatusCode())
		return false
	}
	return true
}
