/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func ensureTargetConnection() error {
	// If the targetConfig.accesstoken is not set or testAccesToken returns false
	if !testAccessToken() {
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
	log.Debugln("Authenticating vRA")
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
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
	log.Debugln("Authenticating vRA Cloud")
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetBody(AuthenticationRequestCloud{target.apitoken}).
		SetResult(&AuthenticationResponseCloud{}).
		SetError(&AuthenticationError{}).
		Post("https://" + target.server + "/iaas/api/login")
	if queryResponse.IsError() {
		log.Debugln("Authentication failed!", queryResponse.RawResponse)
		return "", errors.New(queryResponse.Error().(*AuthenticationError).ServerMessage)
	}
	return queryResponse.Result().(*AuthenticationResponseCloud).Token, err
}

func testAccessToken() bool {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetHeader("Accept", "application/json").
		SetAuthToken(targetConfig.accesstoken).
		SetResult(&UserPreferences{}).
		SetError(&CodeStreamException{}).
		Get("https://" + targetConfig.server + "/pipeline/api/user-preferences")
	if err != nil {
		log.Warnln(err)
		return false
	}
	// log.Debugln(queryResponse.RawResponse)
	if queryResponse.StatusCode() == 401 {
		log.Debugln("Access Token Expired")
		return false
	}
	log.Debugln("Access Token OK (Username:", queryResponse.Result().(*UserPreferences).UserName, ")")
	return true
}

func exportYaml(name, project, path, object string) error {
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
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		SetAuthToken(targetConfig.accesstoken).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		SetError(&CodeStreamException{}).
		Get("https://" + targetConfig.server + "/pipeline/api/export")
	log.Debugln(queryResponse.Request.RawRequest.URL)

	if queryResponse.IsError() {
		return errors.New(queryResponse.Status())
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
	queryResponse, _ := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Content-Type", "application/x-yaml").
		SetBody(yamlPayload).
		SetAuthToken(targetConfig.accesstoken).
		SetError(&CodeStreamException{}).
		Post("https://" + targetConfig.server + "/pipeline/api/import")
	log.Debugln(queryResponse.Request.RawRequest.URL)
	if queryResponse.IsError() {
		return queryResponse.Error().(error)
	}
	var importResponse CodeStreamPipelineImportResponse
	if err = yaml.Unmarshal(queryResponse.Body(), &importResponse); err != nil {
		return err
	}

	if importResponse.Status != "CREATED" && action == "create" {
		return errors.New(importResponse.Status + " - " + importResponse.StatusMessage)
	}
	if importResponse.Status != "UPDATED" && action == "apply" {
		return errors.New(importResponse.Status + " - " + importResponse.StatusMessage)
	}
	return nil
}
