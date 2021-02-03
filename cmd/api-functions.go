package cmd

import (
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
		SetResult(&ExecutionsList{}).
		SetAuthToken(apiKey).
		Get("https://" + server + "/pipeline/api/executions")
	if response.IsError() {
		fmt.Println("GET Executions failed", err)
		os.Exit(1)
	}

	for _, value := range response.Result().(*ExecutionsList).Documents {
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
