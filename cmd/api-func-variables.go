/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func getVariable(id, name, project, exportPath string) ([]*CodeStreamVariableResponse, error) {
	var arrVariables []*CodeStreamVariableResponse
	//var qParams = make(map[string]string)
	client := resty.New()

	// Get by ID
	if id != "" {
		v, e := getVariableByID(id)
		arrVariables = append(arrVariables, v)
		return arrVariables, e
	}
	if name != "" && project != "" {
		qParams["$filter"] = "((name eq '" + name + "') and (project eq '" + project + "'))"
	} else {
		// Get by name
		if name != "" {
			qParams["$filter"] = "(name eq '" + name + "')"
		}
		// Get by project
		if project != "" {
			qParams["$filter"] = "(project eq '" + project + "')"
		}
	}
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/pipeline/api/variables")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	log.Println(queryResponse.Request.URL)

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamVariableResponse{}
		mapstructure.Decode(value, &c)
		arrVariables = append(arrVariables, &c)
		if exportPath != "" {
			exportVariable(c, exportPath)
		}
	}
	return arrVariables, err
}

// getVariableByID - get Code Stream Variable by ID
func getVariableByID(id string) (*CodeStreamVariableResponse, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		log.Println("GET Variable failed", err)
	}
	return queryResponse.Result().(*CodeStreamVariableResponse), err
}

// createVariable - Create a new Code Stream Variable
func createVariable(name string, description string, variableType string, project string, value string) (*CodeStreamVariableResponse, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetBody(
			CodeStreamVariableRequest{
				Project:     project,
				Kind:        "VARIABLE",
				Name:        name,
				Description: description,
				Type:        variableType,
				Value:       value,
			}).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetError(&CodeStreamException{}).
		SetAuthToken(targetConfig.accesstoken).
		Post("https://" + targetConfig.server + "/pipeline/api/variables")
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	}
	return queryResponse.Result().(*CodeStreamVariableResponse), err
}

// updateVariable - Create a new Code Stream Variable
func updateVariable(id string, name string, description string, typename string, value string) (*CodeStreamVariableResponse, error) {
	variable, _ := getVariableByID(id)
	if name != "" {
		variable.Name = name
	}
	if description != "" {
		variable.Description = description
	}
	if typename != "" {
		variable.Type = typename
	}
	if value != "" {
		variable.Value = value
	}
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetBody(variable).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetError(&CodeStreamException{}).
		SetAuthToken(targetConfig.accesstoken).
		Put("https://" + targetConfig.server + "/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	}
	return queryResponse.Result().(*CodeStreamVariableResponse), err
}

// deleteVariable - Delete a Code Stream Variable
func deleteVariable(id string) (*CodeStreamVariableResponse, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetAuthToken(targetConfig.accesstoken).
		Delete("https://" + targetConfig.server + "/pipeline/api/variables/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*CodeStreamVariableResponse), err
}

// exportVariable - Export a variable to YAML
func exportVariable(variable interface{}, exportPath string) {
	var exportFile string
	// variable will be a CodeStreamVariableResponse, so lets remap to CodeStreamVariableRequest
	c := CodeStreamVariableRequest{}
	mapstructure.Decode(variable, &c)
	yaml, err := yaml.Marshal(c)
	if err != nil {
		log.Println("Unable to export variable ", c.Name)
	}

	if filepath.Ext(exportPath) != ".yaml" {
		exportFile = filepath.Join(exportPath, "variables.yaml")
	} else {
		exportFile = exportPath
	}

	file, err := os.OpenFile(exportFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	file.WriteString("---\n" + string(yaml))
}

// importVariables - Import variables from the filePath
func importVariables(filePath string) []CodeStreamVariableRequest {
	var returnVariables []CodeStreamVariableRequest
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(yamlFile)
	decoder := yaml.NewDecoder(reader)
	var request CodeStreamVariableRequest
	for decoder.Decode(&request) == nil {
		returnVariables = append(returnVariables, request)
	}
	return returnVariables
}
