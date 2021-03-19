/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func getProject(id, name, project string) ([]*CodeStreamProject, error) {
	var projects []*CodeStreamProject
	client := resty.New()

	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamProjectList{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/project-service/api/projects")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}

	log.Println(queryResponse.Request.URL)

	for _, value := range queryResponse.Result().(*CodeStreamProjectList).Content {
		c := CodeStreamProject{}
		mapstructure.Decode(value, &c)
		projects = append(projects, &c)
	}
	return projects, err
}
