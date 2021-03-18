/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// PrettyPrint prints interfaces
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// PrintTable prints an array of objects with table headers
func PrintTable(objects []interface{}, headers []string) {
	// Print result table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, object := range objects {
		t := reflect.TypeOf(object)
		fmt.Println(t)
		// var o t
		// mapstructure.Decode(object, &t)
		// var values []string
		// for _, header := range headers {
		// 	append(values, object.(*t).header)
		// }
		// 	table.Append([]string{c.ID, c.Name, c.Project})
	}
	table.Render()
}

func getYamlFilePaths(importPath string) []string {
	var yamlFiles []string
	// Read importPath
	stat, err := os.Stat(importPath)
	if err == nil && stat.IsDir() {
		// log.Debugln("importPath is a directory")
		files, err := ioutil.ReadDir(importPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if strings.Contains(f.Name(), ".yaml") || strings.Contains(f.Name(), ".yml") {
				yamlFiles = append(yamlFiles, filepath.Join(importPath, f.Name()))
			}
		}
	} else {
		// log.Debugln("importPath is a file")
		yamlFiles = append(yamlFiles, importPath)
	}
	return yamlFiles
}
