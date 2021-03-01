/*
Copyright © 2021 Sam McGeown <smcgeown@vmware.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getEndpointCmd represents the endpoint command
var getEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Get Code Stream Endpoint Configurations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()

		response, err := getEndpoint(id, name, project, export, exportPath)
		if err != nil {
			fmt.Print("Unable to get endpoints: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			fmt.Println("No results found")
		} else if resultCount == 1 {
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Project", "Type", "Description"})
			for _, c := range response {
				table.Append([]string{c.Name, c.Project, c.Type, c.Description})
			}
			table.Render()
		}

	},
}

func init() {
	getCmd.AddCommand(getEndpointCmd)
	getEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Get Endpoint by Name")
	getEndpointCmd.Flags().StringVarP(&id, "id", "i", "", "Get Endpoint by ID")
	getEndpointCmd.Flags().StringVarP(&project, "project", "p", "", "Filter Endpoint by Project")
	getEndpointCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getEndpointCmd.Flags().BoolVarP(&export, "export", "e", false, "Export Endpoint")

}
