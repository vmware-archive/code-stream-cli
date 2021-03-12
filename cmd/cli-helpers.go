/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
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
