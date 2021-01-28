package cmd

import (
	"fmt"
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
		fmt.Print("Unable to get an access token: ", response.Error())
	}
	return response.Result().(*AuthenticationResponse).AccessToken, err
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
	fmt.Println(status)

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
		fmt.Println("GET request failed", err)
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
		fmt.Println("GET request failed", err)
	}
	return response.Result().(*CodestreamAPIExecutions), err
}
