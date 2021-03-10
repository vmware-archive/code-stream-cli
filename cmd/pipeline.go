package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var state string
var exportPath string
var importPath string
var export bool
var form bool

// getPipelineCmd represents the pipeline command
var getPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Get Pipelines",
	Long: `Get Code Stream Pipelines by ID, name or status
# List all executions
cs-cli get execution
# View an execution by ID
cs-cli get execution --id 9cc5aedc-db48-4c02-a5e4-086de3160dc0
# View executions of a specific pipeline
get execution --name vra-authenticateUser
# View executions by status
cs-cli get execution --status Failed`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		response, err := getPipelines(id, name, project, export, exportPath)
		if err != nil {
			log.Println("Unable to get Code Stream Pipelines: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Println("No results found")
		} else if resultCount == 1 {
			if form {
				// Get the input form
				var inputs = response[0].Input
				PrettyPrint(inputs)
			} else {
				// Print the single result
				PrettyPrint(response[0])
			}
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project})
			}
			table.Render()
		}
	},
}

// updatePipelineCmd represents the pipeline update command
var updatePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Update a Pipeline",
	Long: `Update a Pipeline
	Enable/Disable/Release:
	cs-cli update pipeline --id d0185f04-2e87-4f3c-b6d7-ee58abba3e92 --state enabled/disabled/released
	Update from YAML
	cs-cli update pipeline --importPath "/Users/sammcgeown/Desktop/pipelines/SSH Exports.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if state != "" {
			switch strings.ToUpper(state) {
			case "ENABLED", "DISABLED", "RELEASED":
				// Valid states
				return nil
			}
			return errors.New("--state is not valid, must be ENABLED, DISABLED or RELEASED")
		}
		if export {

		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if state != "" {
			response, err := patchPipeline(id, `{"state":"`+state+`"}`)
			if err != nil {
				log.Println("Unable to update Code Stream Pipeline: ", err)
			}
			log.Println("Setting pipeline " + response.Name + " to " + state)
		}

		if importPath != "" {
			err := importYaml(importPath, "apply")
			if err != nil {
				log.Fatalln("Failed to update Pipeline", err)
			}
			log.Println("Imported successfully, pipeline updated.")
		}
	},
}

// createPipelineCmd represents the pipeline create command
var createPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Create a Pipeline",
	Long: `Create a Pipeline by importing a YAML specification.
	
	Create from YAML
	  cs-cli create pipeline --importPath "/Users/sammcgeown/Desktop/pipelines/SSH Exports.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if importPath != "" {
			err := importYaml(importPath, "create")
			if err != nil {
				log.Fatalln("Failed to import Pipeline", err)
			}
			fmt.Println("Imported successfully, Pipeline created.")
		}
	},
}

// deletePipelineCmd represents the delete pipeline command
var deletePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Delete a Pipeline",
	Long: `Delete a Pipeline with a specific ID
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		response, err := deletePipeline(id)
		if err != nil {
			log.Fatalln("Delete Pipeline failed:", err)
		}
		log.Println("Pipeline with id " + response.ID + " deleted")

	},
}

func init() {
	// Get
	getCmd.AddCommand(getPipelineCmd)
	getPipelineCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to list executions for")
	getPipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	getPipelineCmd.Flags().StringVarP(&project, "project", "p", "", "List pipeline in project")
	getPipelineCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getPipelineCmd.Flags().BoolVarP(&export, "export", "e", false, "Export pipeline")
	getPipelineCmd.Flags().BoolVarP(&form, "form", "f", false, "Get pipeline inputs")

	// Create
	createCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().StringVarP(&importPath, "importPath", "", "", "YAML configuration file to import")
	// createPipelineCmd.Flags().StringVarP(&project, "project", "p", "", "Manually specify the Project in which to create the Pipeline (overrides YAML)")
	createPipelineCmd.MarkFlagRequired("importPath")
	// Update
	updateCmd.AddCommand(updatePipelineCmd)
	updatePipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	updatePipelineCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	updatePipelineCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the pipeline (ENABLED|DISABLED|RELEASED")
	// Delete
	deleteCmd.AddCommand(deletePipelineCmd)
	deletePipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Pipeline to delete")
	deletePipelineCmd.MarkFlagRequired("id")

}
