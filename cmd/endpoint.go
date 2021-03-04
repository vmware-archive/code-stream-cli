package cmd

import (
	"log"
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

		response, err := getEndpoint(id, name, project, typename, export, exportPath)
		if err != nil {
			log.Println("Unable to get endpoints: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Println("No results found")
		} else if resultCount == 1 {
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Project", "Type", "Description"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description})
			}
			table.Render()
		}

	},
}

// createEndpointCmd represents the endpoint create command
var createEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Create an Endpoint",
	Long: `Create an Endpoint by importing a YAML specification.
	
	Create from YAML
	  cs-cli create endpoint --importPath "/Users/sammcgeown/Desktop/endpoint.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		if importPath != "" {
			err := importYaml(importPath, "create")
			if err != nil {
				log.Fatalln("Failed to import Endpoint", err)
			}
			log.Println("Imported successfully, Endpoint created.")
		}
	},
}

// updateEndpointCmd represents the endpoint update command
var updateEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Update an Endpoint",
	Long: `Update an Endpoint by importing the YAML specification

	Update from YAML
	cs-cli update endpoint --importPath "/Users/sammcgeown/Desktop/updated-endpoint.yaml"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		if importPath != "" {
			err := importYaml(importPath, "update")
			if err != nil {
				log.Fatalln("Failed to import Endpoint", err)
			}
			log.Println("Imported successfully, Endpoint updated.")
		}
	},
}

// deleteEndpointCmd represents the executions command
var deleteEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Delete an Endpoint",
	Long: `Delete an Endpoint with a specific Endpoint ID
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()

		response, err := deleteEndpoint(id)
		if err != nil {
			log.Println("Unable to delete Endpoint: ", err)
		}
		log.Println("Endpoint with id " + response.ID + " deleted")
	},
}

func init() {
	getCmd.AddCommand(getEndpointCmd)
	getEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Get Endpoint by Name")
	getEndpointCmd.Flags().StringVarP(&id, "id", "i", "", "Get Endpoint by ID")
	getEndpointCmd.Flags().StringVarP(&project, "project", "p", "", "Filter Endpoint by Project")
	getEndpointCmd.Flags().StringVarP(&typename, "type", "t", "", "Filter Endpoint by Type")
	getEndpointCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getEndpointCmd.Flags().BoolVarP(&export, "export", "e", false, "Export Endpoint")
	// Create
	createCmd.AddCommand(createEndpointCmd)
	createEndpointCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	createEndpointCmd.MarkFlagRequired("importPath")
	// Update
	updateCmd.AddCommand(updateEndpointCmd)
	updateEndpointCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	updateEndpointCmd.MarkFlagRequired("importPath")
	// Delete
	deleteCmd.AddCommand(deleteEndpointCmd)
	deleteEndpointCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Endpoint to delete")
	deleteEndpointCmd.MarkFlagRequired("id")

}
