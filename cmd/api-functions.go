package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func authenticate(server string, username string, password string, domain string) (string, error) {
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

func testAccessToken() bool {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		Get("https://" + server + "/iaas/api/projects")
	if response.IsError() {
		fmt.Println("Key test failed", err)
		return false
	}
	if response.StatusCode() == 401 {
		//fmt.Println("Token authentication failed: ", response.StatusCode())
		return false
	}
	return true
}

func getExecutions(id string, status string) ([]*CodestreamAPIExecutions, error) {
	var arrExecutions []*CodestreamAPIExecutions
	if id != "" {
		x, err := getExecution("/codestream/api/executions/" + id)
		if err != nil {
			fmt.Print("Error: ", err.Error())
		}
		arrExecutions = append(arrExecutions, x)
		return arrExecutions, err
	}
	client := resty.New()
	var qParams = make(map[string]string)
	qParams["$orderby"] = "_requestTimeInMicros desc"
	if status != "" {
		qParams["$filter"] = "((status eq '" + strings.ToUpper(status) + "')) and _nested eq 'false'"
	}
	response, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(apiKey).
		Get("https://" + server + "/pipeline/api/executions")
	if response.IsError() {
		fmt.Println("GET Executions failed", err)
		os.Exit(1)
	}

	for _, value := range response.Result().(*documentsList).Documents {
		c := CodestreamAPIExecutions{}
		mapstructure.Decode(value, &c)
		arrExecutions = append(arrExecutions, &c)
	}
	return arrExecutions, err
}

func getExecution(executionLink string) (*CodestreamAPIExecutions, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&CodestreamAPIExecutions{}).
		SetAuthToken(apiKey).
		Get("https://" + server + executionLink)
	if response.IsError() {
		fmt.Println("GET Execution failed", err)
	}
	return response.Result().(*CodestreamAPIExecutions), err
}

func deleteExecution(id string) (*CodestreamAPIExecutions, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&CodestreamAPIExecutions{}).
		SetAuthToken(apiKey).
		Delete("https://" + server + "/pipeline/api/executions/" + id)
	if response.IsError() {
		fmt.Println("DELETE Execution failed", err)
	}
	return response.Result().(*CodestreamAPIExecutions), err
}

func getVariable(id, name, project string) ([]*CodeStreamVariableResponse, error) {
	var arrVariables []*CodeStreamVariableResponse
	var qParams = make(map[string]string)
	client := resty.New()

	// Get by ID
	if id != "" {
		v, e := getVariableByID(id)
		arrVariables = append(arrVariables, v)
		return arrVariables, e
	}
	// Get by name
	if name != "" {
		qParams["$filter"] = "(name eq '" + name + "')"
	}
	// Get by project
	if project != "" {
		qParams["$filter"] = "(project eq '" + project + "')"
	}
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(apiKey).
		Get("https://" + server + "/pipeline/api/variables")

	if queryResponse.IsError() {
		fmt.Println("GET Variables failed", err)
		os.Exit(1)
	}

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamVariableResponse{}
		mapstructure.Decode(value, &c)
		arrVariables = append(arrVariables, &c)
	}
	return arrVariables, err
}

// getVariableByID - get Code Stream Variable by ID
func getVariableByID(id string) (*CodeStreamVariableResponse, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetAuthToken(apiKey).
		Get("https://" + server + "/pipeline/api/variables/" + id)
	if response.IsError() {
		fmt.Println("GET Variable failed", err)
	}
	return response.Result().(*CodeStreamVariableResponse), err
}

// createVariable - Create a new Code Stream Variable
func createVariable(name string, description string, variableType string, project string, value string) (*CodeStreamVariableResponse, error) {
	client := resty.New()
	response, err := client.R().
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
		SetAuthToken(apiKey).
		Post("https://" + server + "/pipeline/api/variables")
	if response.IsError() {
		fmt.Println("Create Variable failed", err)
		os.Exit(1)
	}
	return response.Result().(*CodeStreamVariableResponse), err
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
	response, err := client.R().
		SetBody(variable).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetAuthToken(apiKey).
		Put("https://" + server + "/pipeline/api/variables/" + id)
	if response.IsError() {
		fmt.Println("Create Variable failed", err)
		os.Exit(1)
	}
	return response.Result().(*CodeStreamVariableResponse), err
}

// deleteVariable - Delete a Code Stream Variable
func deleteVariable(id string) (*CodeStreamVariableResponse, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamVariableResponse{}).
		SetAuthToken(apiKey).
		Delete("https://" + server + "/pipeline/api/variables/" + id)
	if response.IsError() {
		fmt.Println("Create Variable failed", err)
		os.Exit(1)
	}
	return response.Result().(*CodeStreamVariableResponse), err
}

// PrettyPrint prints interfaces
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func getPipelines(id, name, project string) ([]*CodeStreamPipeline, error) {
	var arrResults []*CodeStreamPipeline
	var qParams = make(map[string]string)
	client := resty.New()

	// Get by ID
	if id != "" {
		v, e := getPipelineByID(id)
		arrResults = append(arrResults, v)
		return arrResults, e
	}
	// Get by name
	if name != "" {
		qParams["$filter"] = "(name eq '" + name + "')"
	}
	// Get by project
	if project != "" {
		qParams["$filter"] = "(project eq '" + project + "')"
	}
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(apiKey).
		Get("https://" + server + "/pipeline/api/pipelines")

	if queryResponse.IsError() {
		fmt.Println("GET Variables failed", err)
		os.Exit(1)
	}

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamPipeline{}
		mapstructure.Decode(value, &c)
		arrResults = append(arrResults, &c)
	}
	return arrResults, err
}
func exportPipeline(name, project string) {
	var qParams = make(map[string]string)
	qParams["pipelines"] = name
	qParams["project"] = project
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/x-yaml;charset=UTF-8").
		SetAuthToken(apiKey).
		SetOutput(name + ".yaml").
		Get("https://" + server + "/pipeline/api/export")

	if queryResponse.IsError() {
		fmt.Println("Exit pipeline failed", err)
		os.Exit(1)
	}
}

// getPipelineByID - get Code Stream Pipeline by ID
func getPipelineByID(id string) (*CodeStreamPipeline, error) {
	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(apiKey).
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
		SetAuthToken(apiKey).
		Patch("https://" + server + "/pipeline/api/pipelines/" + id)
	if response.IsError() {
		fmt.Println("GET Pipeline failed", response.StatusCode())
		return nil, err
	}
	return response.Result().(*CodeStreamPipeline), nil
}
