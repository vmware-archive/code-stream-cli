package cmd

import (
	"fmt"
	"os"

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
		if export {
			exportYaml(c.Name, c.Project, exportPath, "piplines")
			endpoints = append(endpoints, &c)
		} else {
			endpoints = append(endpoints, &c)
		}

	}
	return endpoints, err
}
