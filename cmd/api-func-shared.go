/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func ensureTargetConnection() error {
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
			return authError
		}
		if viper.ConfigFileUsed() != "" {
			viper.Set("target."+currentTargetName+".accesstoken", targetConfig.accesstoken)
			viper.WriteConfig()
		}
	}
	return nil
}

func authenticateOnPrem(target config) (string, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetBody(AuthenticationRequest{target.username, target.password, target.domain}).
		SetResult(&AuthenticationResponse{}).
		SetError(&AuthenticationError{}).
		Post("https://" + target.server + "/csp/gateway/am/idp/auth/login?access_token")
	if queryResponse.IsError() {
		return "", errors.New(queryResponse.Error().(*AuthenticationError).ServerMessage)
	}
	return queryResponse.Result().(*AuthenticationResponse).AccessToken, err
}
func authenticateCloud(target config) (string, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetBody(AuthenticationRequestCloud{target.apitoken}).
		SetResult(&AuthenticationResponseCloud{}).
		SetError(&AuthenticationError{}).
		Post("https://" + target.server + "/iaas/api/login")
	if queryResponse.IsError() {
		return "", errors.New(queryResponse.Error().(*AuthenticationError).ServerMessage)
	}
	return queryResponse.Result().(*AuthenticationResponseCloud).Token, err
}

func testAccessToken() bool {
	client := resty.New()
	queryResponse, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/iaas/api/projects")
	if err != nil {
		return false
	}
	if queryResponse.StatusCode() == 401 {
		return false
	}
	return true
}

func exportYaml(name, project, path, object string) error {
	var exportPath string
	qParams[object] = name
	qParams["project"] = project
	if path != "" {
		exportPath = path
	} else {
		exportPath, _ = os.Getwd()
	}
	client := resty.New()
	queryResponse, _ := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		SetAuthToken(targetConfig.accesstoken).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		Get("https://" + targetConfig.server + "/pipeline/api/export")

	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}
	return nil
}

// importYaml import a yaml pipeline or endpoint
func importYaml(yamlPath, action string) error {
	qParams["action"] = action
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	yamlPayload := string(yamlBytes)
	client := resty.New()
	queryResponse, _ := client.R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/x-yaml").
		SetBody(yamlPayload).
		SetAuthToken(targetConfig.accesstoken).
		Post("https://" + targetConfig.server + "/pipeline/api/import")
	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}
	return nil
}
