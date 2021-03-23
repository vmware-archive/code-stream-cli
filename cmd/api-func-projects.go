/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"crypto/tls"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func getProject(id, name string) ([]*CodeStreamProject, error) {
	var projects []*CodeStreamProject
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

	queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
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
