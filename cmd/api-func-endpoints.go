package cmd

import (
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func getEndpoint(id, name, project, endpointtype string, export bool, exportPath string) ([]*CodeStreamEndpoint, error) {
	var endpoints []*CodeStreamEndpoint
	var qParams = make(map[string]string)
	qParams["expand"] = "true"
	client := resty.New()

	var filters []string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if endpointtype != "" {
		filters = append(filters, "(type eq '"+endpointtype+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	}

	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/pipeline/api/endpoints")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamEndpoint{}
		mapstructure.Decode(value, &c)
		if export {
			exportYaml(c.Name, c.Project, exportPath, "endpoints")
			endpoints = append(endpoints, &c)
		} else {
			endpoints = append(endpoints, &c)
		}

	}
	return endpoints, err
}

func deleteEndpoint(id string) (*CodeStreamEndpoint, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamEndpoint{}).
		SetAuthToken(targetConfig.accesstoken).
		Delete("https://" + targetConfig.server + "/pipeline/api/endpoints/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*CodeStreamEndpoint), err
}
