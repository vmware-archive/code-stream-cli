/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func getPipelines(id string, name string, project string, exportPath string) ([]*CodeStreamPipeline, error) {
	var arrResults []*CodeStreamPipeline
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
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	}
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(targetConfig.accesstoken).
		SetError(&CodeStreamException{}).
		Get("https://" + targetConfig.server + "/pipeline/api/pipelines")

	log.Debugln(queryResponse.Request.RawRequest.URL)
	// log.Debugln(queryResponse.String())

	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)

	}
	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamPipeline{}
		mapstructure.Decode(value, &c)
		if exportPath != "" {
			if err := exportYaml(c.Name, c.Project, exportPath, "pipelines"); err != nil {
				log.Warnln(err)
			}
			arrResults = append(arrResults, &c)
		} else {
			arrResults = append(arrResults, &c)
		}
	}
	return arrResults, err
}

// patchPipeline - Patch Code Stream Pipeline by ID
func patchPipeline(id string, payload string) (*CodeStreamPipeline, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(targetConfig.accesstoken).
		Patch("https://" + targetConfig.server + "/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*CodeStreamPipeline), err
}

func deletePipeline(id string) (*CodeStreamPipeline, error) {
	client := resty.New()
	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(targetConfig.accesstoken).
		SetError(&CodeStreamException{}).
		Delete("https://" + targetConfig.server + "/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, errors.New(queryResponse.Error().(*CodeStreamException).Message)
	}
	return queryResponse.Result().(*CodeStreamPipeline), err
}

func deletePipelineInProject(project string) ([]*CodeStreamPipeline, error) {
	var deletedPipes []*CodeStreamPipeline
	pipelines, err := getPipelines("", "", project, "")
	if err != nil {
		return nil, err
	}
	confirm := askForConfirmation("This will attempt to delete " + fmt.Sprint(len(pipelines)) + " Pipelines in " + project + ", are you sure?")
	if confirm {
		for _, pipeline := range pipelines {
			deletedPipe, err := deletePipeline(pipeline.ID)
			if err != nil {
				log.Warnln("Unable to delete "+pipeline.Name, err)
			}
			deletedPipes = append(deletedPipes, deletedPipe)
		}
		return deletedPipes, nil
	} else {
		return nil, errors.New("user declined")
	}
}
