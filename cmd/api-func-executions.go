package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func getExecutions(id string, status string, name string, nested bool) ([]*CodestreamAPIExecutions, error) {
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
		qParams["$filter"] = "((status eq '" + strings.ToUpper(status) + "') and (_nested eq '" + strconv.FormatBool(nested) + "'))"
	}
	if name != "" {
		qParams["$filter"] = "((name eq '" + name + "') and (_nested eq '" + strconv.FormatBool(nested) + "'))"
	}
	response, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(accessToken).
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
		SetAuthToken(accessToken).
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
		SetAuthToken(accessToken).
		Delete("https://" + server + "/pipeline/api/executions/" + id)
	if response.IsError() {
		fmt.Println("DELETE Execution failed", err)
	}
	return response.Result().(*CodestreamAPIExecutions), err
}

func createExecution(id string, inputs string, comment string) (*CodeStreamCreateExecutionResponse, error) {
	// Convert JSON string to byte array
	var inputBytes = []byte(inputs)
	// Unmarshal inputs using a generic interface
	var inputsInterface interface{}
	err := json.Unmarshal(inputBytes, &inputsInterface)
	if err != nil {
		return nil, err
	}
	// Create CodeStreamCreateExecutionRequest struct
	var execution CodeStreamCreateExecutionRequest
	execution.Comments = comment
	execution.Input = inputsInterface
	//Marshal struct to JSON []byte
	executionBytes, err := json.Marshal(execution)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(executionBytes).
		SetResult(&CodeStreamCreateExecutionResponse{}).
		SetAuthToken(accessToken).
		Post("https://" + server + "/pipeline/api/pipelines/" + id + "/executions")
	fmt.Println(response.StatusCode())
	if response.IsError() {
		return nil, response.Error().(error)
	}
	return response.Result().(*CodeStreamCreateExecutionResponse), nil
}
