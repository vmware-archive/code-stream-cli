package cmd

import (
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func getCustomIntegration(id, name string) ([]*CodeStreamCustomIntegration, error) {
	var arrCustomIntegrations []*CodeStreamCustomIntegration
	client := resty.New()

	var filters []string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	}

	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/pipeline/api/custom-integrations")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamCustomIntegration{}
		mapstructure.Decode(value, &c)
		arrCustomIntegrations = append(arrCustomIntegrations, &c)
	}
	return arrCustomIntegrations, err
}

// // createCustomIntegration - Create a new Code Stream CustomIntegration
// func createCustomIntegration(name string, description string, variableType string, project string, value string) (*CodeStreamCustomIntegrationResponse, error) {
// 	client := resty.New()
// 	response, err := client.R().
// 		SetBody(
// 			CodeStreamCustomIntegrationRequest{
// 				Project:     project,
// 				Kind:        "VARIABLE",
// 				Name:        name,
// 				Description: description,
// 				Type:        variableType,
// 				Value:       value,
// 			}).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&CodeStreamCustomIntegrationResponse{}).
// 		SetError(&CodeStreamException{}).
// 		SetAuthToken(targetConfig.accesstoken).
// 		Post("https://" + targetConfig.server + "/pipeline/api/variables")
// 	if response.IsError() {
// 		return nil, errors.New(response.Error().(*CodeStreamException).Message)
// 	}
// 	return response.Result().(*CodeStreamCustomIntegrationResponse), err
// }

// // updateCustomIntegration - Create a new Code Stream CustomIntegration
// func updateCustomIntegration(id string, name string, description string, typename string, value string) (*CodeStreamCustomIntegrationResponse, error) {
// 	variable, _ := getCustomIntegrationByID(id)
// 	if name != "" {
// 		variable.Name = name
// 	}
// 	if description != "" {
// 		variable.Description = description
// 	}
// 	if typename != "" {
// 		variable.Type = typename
// 	}
// 	if value != "" {
// 		variable.Value = value
// 	}
// 	client := resty.New()
// 	response, err := client.R().
// 		SetBody(variable).
// 		SetHeader("Accept", "application/json").
// 		SetResult(&CodeStreamCustomIntegrationResponse{}).
// 		SetError(&CodeStreamException{}).
// 		SetAuthToken(targetConfig.accesstoken).
// 		Put("https://" + targetConfig.server + "/pipeline/api/variables/" + id)
// 	if response.IsError() {
// 		return nil, errors.New(response.Error().(*CodeStreamException).Message)
// 	}
// 	return response.Result().(*CodeStreamCustomIntegrationResponse), err
// }

// // deleteCustomIntegration - Delete a Code Stream CustomIntegration
// func deleteCustomIntegration(id string) (*CodeStreamCustomIntegrationResponse, error) {
// 	client := resty.New()
// 	response, err := client.R().
// 		SetHeader("Accept", "application/json").
// 		SetResult(&CodeStreamCustomIntegrationResponse{}).
// 		SetAuthToken(targetConfig.accesstoken).
// 		Delete("https://" + targetConfig.server + "/pipeline/api/variables/" + id)
// 	if response.IsError() {
// 		log.Println("Create CustomIntegration failed", err)
// 		os.Exit(1)
// 	}
// 	return response.Result().(*CodeStreamCustomIntegrationResponse), err
// }

// // exportCustomIntegration - Export a variable to YAML
// func exportCustomIntegration(variable interface{}, exportFile string) {
// 	// variable will be a CodeStreamCustomIntegrationResponse, so lets remap to CodeStreamCustomIntegrationRequest
// 	c := CodeStreamCustomIntegrationRequest{}
// 	mapstructure.Decode(variable, &c)
// 	yaml, err := yaml.Marshal(c)
// 	if err != nil {
// 		log.Println("Unable to export variable ", c.Name)
// 	}
// 	if exportFile == "" {
// 		exportFile = "variables.yaml"
// 	}
// 	file, err := os.OpenFile(exportFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer file.Close()
// 	file.WriteString("---\n" + string(yaml))
// }

// // importCustomIntegrations - Import variables from the filePath
// func importCustomIntegrations(filePath string) []CodeStreamCustomIntegrationRequest {
// 	var returnCustomIntegrations []CodeStreamCustomIntegrationRequest
// 	filename, _ := filepath.Abs(filePath)
// 	yamlFile, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	reader := bytes.NewReader(yamlFile)
// 	decoder := yaml.NewDecoder(reader)
// 	var request CodeStreamCustomIntegrationRequest
// 	for decoder.Decode(&request) == nil {
// 		returnCustomIntegrations = append(returnCustomIntegrations, request)
// 	}
// 	return returnCustomIntegrations
// }
