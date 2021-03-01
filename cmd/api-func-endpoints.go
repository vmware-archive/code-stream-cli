package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func getEndpoint(id, name, project string) ([]*CodeStreamEndpoint, error) {
	var endpoints []*CodeStreamEndpoint
	var qParams = make(map[string]string)
	client := resty.New()

	// Get by ID
	// if id != "" {
	// 	v, e := getVariableByID(id)
	// 	endpoints = append(endpoints, v)
	// 	return endpoints, e
	// }
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
		SetAuthToken(accessToken).
		Get("https://" + server + "/pipeline/api/endpoints")

	if queryResponse.IsError() {
		fmt.Println("GET Variables failed", err)
		os.Exit(1)
	}

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamEndpoint{}
		mapstructure.Decode(value, &c)
		endpoints = append(endpoints, &c)
	}
	return endpoints, err
}

func exportEndpoint(name, project, path string) {
	var exportPath string
	var qParams = make(map[string]string)
	qParams["pipelines"] = name
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
		SetAuthToken(accessToken).
		SetOutput(filepath.Join(exportPath, name+".yaml")).
		Get("https://" + server + "/pipeline/api/export")

	if queryResponse.IsError() {
		fmt.Println("Export pipeline failed", err)
		os.Exit(1)
	}
}
