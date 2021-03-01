package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func getPipelines(id string, name string, project string, export bool, exportPath string) ([]*CodeStreamPipeline, error) {
	var arrResults []*CodeStreamPipeline
	var qParams = make(map[string]string)
	client := resty.New()
	// Get by ID
	if id != "" {
		v, e := getPipelineByID(id)
		arrResults = append(arrResults, v)
		return arrResults, e
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
		SetAuthToken(accessToken).
		Get("https://" + server + "/pipeline/api/pipelines")

	if queryResponse.IsError() {
		fmt.Println("GET Variables failed", err)
		os.Exit(1)
	}
	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamPipeline{}
		mapstructure.Decode(value, &c)
		if export {
			exportYaml(c.Name, c.Project, exportPath, "piplines")
			arrResults = append(arrResults, &c)
		} else {
			arrResults = append(arrResults, &c)
		}
	}
	return arrResults, err
}

// getPipelineByID - get Code Stream Pipeline by ID
func getPipelineByID(id string) (*CodeStreamPipeline, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(accessToken).
		Get("https://" + server + "/pipeline/api/pipelines/" + id)
	if response.IsError() {
		fmt.Println("GET Pipeline failed", err)
	}
	return response.Result().(*CodeStreamPipeline), err
}

// patchPipeline - Patch Code Stream Pipeline by ID
func patchPipeline(id string, payload string) (*CodeStreamPipeline, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(accessToken).
		Patch("https://" + server + "/pipeline/api/pipelines/" + id)
	if response.IsError() {
		fmt.Println("GET Pipeline failed", response.StatusCode())
		return nil, err
	}
	return response.Result().(*CodeStreamPipeline), nil
}
