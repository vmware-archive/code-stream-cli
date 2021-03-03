/*
Package cmd Copyright © 2021 Sam McGeown <smcgeown@vmware.com>

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
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getCustomIntegrationCmd represents the customintegration command
var getCustomIntegrationCmd = &cobra.Command{
	Use:   "customintegration",
	Short: "Get vRealize Code Stream Custom Integrations",
	Long: `Get vRealize Code Stream Custom Integrations by name, project or by id - e.g:

Get by ID
	cs-cli get customintegration --id 6b7936d3-a19d-4298-897a-65e9dc6620c8
	
Get by Name
	cs-cli get customintegration --name my-customintegration
	
Get by Project
	cs-cli get customintegration --project production`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		response, err := getCustomIntegration(id, name)
		if err != nil {
			log.Println("Unable to get Code Stream CustomIntegrations: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Println("No results found")
		} else if resultCount == 1 {
			// Print the single result
			if export {
				//exportCustomIntegration(response[0], exportFile)
			}
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Status", "Description"})
			for _, c := range response {
				if export {
					//exportCustomIntegration(c, exportFile)
				}
				table.Append([]string{c.ID, c.Name, c.Status, c.Description})
			}
			table.Render()
		}
	},
}

// // getCustomIntegrationCmd represents the customintegration command
// var createCustomIntegrationCmd = &cobra.Command{
// 	Use:   "customintegration",
// 	Short: "A brief description of your command",
// 	Long:  ``,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		ensureTargetConnection()

// 		if importFile != "" { // If we are importing a file
// 			customintegrations := importCustomIntegrations(importFile)
// 			for _, value := range customintegrations {
// 				if project != "" { // If the project is specified update the object
// 					value.Project = project
// 				}
// 				createResponse, err := createCustomIntegration(value.Name, value.Description, value.Type, value.Project, value.Value)
// 				if err != nil {
// 					log.Println("Unable to create Code Stream CustomIntegration: ", err)
// 				} else {
// 					log.Println("Created customintegration", createResponse.Name, "in", createResponse.Project)
// 				}
// 			}
// 		} else {
// 			createResponse, err := createCustomIntegration(name, description, typename, project, value)
// 			if err != nil {
// 				log.Println("Unable to create Code Stream CustomIntegration: ", err)
// 			}
// 			PrettyPrint(createResponse)
// 		}
// 	},
// }

// // updateCustomIntegrationCmd represents the customintegration command
// var updateCustomIntegrationCmd = &cobra.Command{
// 	Use:   "customintegration",
// 	Short: "A brief description of your command",
// 	Long:  ``,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		ensureTargetConnection()
// 		if importFile != "" { // If we are importing a file
// 			customintegrations := importCustomIntegrations(importFile)
// 			for _, value := range customintegrations {
// 				exisitingCustomIntegration, err := getCustomIntegration("", value.Name, value.Project)
// 				if err != nil {
// 					log.Println("Update failed - unable to find existing Code Stream CustomIntegration", value.Name, "in", value.Project)
// 				} else {
// 					_, err := updateCustomIntegration(exisitingCustomIntegration[0].ID, value.Name, value.Description, value.Type, value.Value)
// 					if err != nil {
// 						log.Println("Unable to update Code Stream CustomIntegration: ", err)
// 					} else {
// 						log.Println("Updated customintegration", value.Name)
// 					}
// 				}
// 			}
// 		} else { // Else we are updating using flags
// 			updateResponse, err := updateCustomIntegration(id, name, description, typename, value)
// 			if err != nil {
// 				log.Println("Unable to update Code Stream CustomIntegration: ", err)
// 			}
// 			log.Println("Updated customintegration", updateResponse.Name)
// 		}
// 	},
// }

// // deleteCustomIntegrationCmd represents the executions command
// var deleteCustomIntegrationCmd = &cobra.Command{
// 	Use:   "customintegration",
// 	Short: "A brief description of your command",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:

// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		ensureTargetConnection()

// 		response, err := deleteCustomIntegration(id)
// 		if err != nil {
// 			log.Println("Unable to delete customintegration: ", err)
// 		}
// 		log.Println("CustomIntegration with id " + response.ID + " deleted")
// 	},
// }

func init() {
	// Get CustomIntegration
	getCmd.AddCommand(getCustomIntegrationCmd)
	getCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "List customintegration with name")
	getCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "List customintegrations by id")
	getCustomIntegrationCmd.Flags().StringVarP(&exportFile, "exportFile", "", "", "Path to export objects - relative or absolute location")
	getCustomIntegrationCmd.Flags().BoolVarP(&export, "export", "e", false, "Export customintegrations, uses ./customintegrations.yaml or the file specified by --exportFile")
	// // Create CustomIntegration
	// createCmd.AddCommand(createCustomIntegrationCmd)
	// createCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the customintegration to create")
	// createCustomIntegrationCmd.Flags().StringVarP(&typename, "type", "t", "", "The type of the customintegration to create REGULAR|SECRET|RESTRICTED")
	// createCustomIntegrationCmd.Flags().StringVarP(&project, "project", "p", "", "The project in which to create the customintegration")
	// createCustomIntegrationCmd.Flags().StringVarP(&value, "value", "v", "", "The value of the customintegration to create")
	// createCustomIntegrationCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the customintegration to create")
	// createCustomIntegrationCmd.Flags().StringVarP(&importFile, "importfile", "i", "", "Path to a YAML file with the customintegrations to import")

	// // Update CustomIntegration
	// updateCmd.AddCommand(updateCustomIntegrationCmd)
	// updateCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the customintegration to update")
	// updateCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "Update the name of the customintegration")
	// updateCustomIntegrationCmd.Flags().StringVarP(&typename, "type", "t", "", "Update the type of the customintegration REGULAR|SECRET|RESTRICTED")
	// updateCustomIntegrationCmd.Flags().StringVarP(&value, "value", "v", "", "Update the value of the customintegration ")
	// updateCustomIntegrationCmd.Flags().StringVarP(&description, "description", "d", "", "Update the description of the customintegration")
	// updateCustomIntegrationCmd.Flags().StringVarP(&importFile, "importfile", "", "", "Path to a YAML file with the customintegrations to import")
	// //updateCustomIntegrationCmd.MarkFlagRequired("id")
	// // Delete CustomIntegration
	// deleteCmd.AddCommand(deleteCustomIntegrationCmd)
	// deleteCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "Delete customintegration by id")
	// deleteCustomIntegrationCmd.MarkFlagRequired("id")
}
