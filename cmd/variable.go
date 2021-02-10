/*
Package cmd Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

var id string
var name string
var project string
var typename string
var value string
var description string

// getVariableCmd represents the variable command
var getVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Get vRealize Code Stream Variables",
	Long: `Get vRealize Code Stream Variables by name, project or by id - e.g:

Get by ID
	cs-cli get variable --id 6b7936d3-a19d-4298-897a-65e9dc6620c8
	
Get by Name
	cs-cli get variable --name my-variable
	
Get by Project
	cs-cli get variable --project production`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := getVariable(id, name, project)
		if err != nil {
			fmt.Print("Unable to get Code Stream Variables: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			fmt.Println("No results found")
		} else if resultCount == 1 {
			// Print the single result
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Type", "Description"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description})
			}
			table.Render()
		}
	},
}

// getVariableCmd represents the variable command
var createVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		createResponse, err := createVariable(name, description, typename, project, value)
		if err != nil {
			fmt.Print("Unable to create Code Stream Variable: ", err)
		}

		PrettyPrint(createResponse)
	},
}

// updateVariableCmd represents the variable command
var updateVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		updateResponse, err := updateVariable(id, name, description, typename, value)
		if err != nil {
			fmt.Print("Unable to update Code Stream Variable: ", err)
		}

		PrettyPrint(updateResponse)
	},
}

// deleteVariableCmd represents the executions command
var deleteVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := deleteVariable(i)
		if err != nil {
			fmt.Print("Unable to delete variable: ", err)
		}
		fmt.Println("Variable with id " + response.ID + " deleted")
	},
}

func init() {
	// Get Variable
	getCmd.AddCommand(getVariableCmd)
	getVariableCmd.Flags().StringVarP(&name, "name", "n", "", "List variable with name")
	getVariableCmd.Flags().StringVarP(&project, "project", "p", "", "List variables in project")
	getVariableCmd.Flags().StringVarP(&id, "id", "i", "", "List variables by id")
	// Create Variable
	createCmd.AddCommand(createVariableCmd)
	createVariableCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the variable to create")
	createVariableCmd.Flags().StringVarP(&typename, "type", "t", "", "The type of the variable to create REGULAR|SECRET|RESTRICTED")
	createVariableCmd.Flags().StringVarP(&project, "project", "p", "", "The project in which to create the variable")
	createVariableCmd.Flags().StringVarP(&value, "value", "v", "", "The value of the variable to create")
	createVariableCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the variable to create")
	// Update Variable
	updateCmd.AddCommand(updateVariableCmd)
	updateVariableCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the variable to update")
	updateVariableCmd.Flags().StringVarP(&name, "name", "n", "", "Update the name of the variable")
	updateVariableCmd.Flags().StringVarP(&typename, "type", "t", "", "Update the type of the variable REGULAR|SECRET|RESTRICTED")
	updateVariableCmd.Flags().StringVarP(&value, "value", "v", "", "Update the value of the variable ")
	updateVariableCmd.Flags().StringVarP(&description, "description", "d", "", "Update the description of the variable")
	updateVariableCmd.MarkFlagRequired("id")
	// Delete Variable
	deleteCmd.AddCommand(deleteVariableCmd)
	deleteVariableCmd.Flags().StringVarP(&id, "id", "i", "", "Delete variable by id")
	deleteVariableCmd.MarkFlagRequired("id")
}
