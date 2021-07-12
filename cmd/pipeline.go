/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var state string
var printForm bool
var dependencies bool

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

		response, err := getPipelines(id, name, project, exportPath)
		if err != nil {
			log.Errorln("Unable to get Code Stream Pipelines: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Warnln("No results found")
		}

		if printJson {
			for _, c := range response {
				PrettyPrint(c)
			}
		} else if printForm {
			// Get the input form
			for _, c := range response {
				PrettyPrint(c.Input)
			}
		} else {
			var endpoints []string
			var variables []string
			var pipelines []string
			var customintegrations []string

			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Description"})
			for _, c := range response {
				if c.Workspace.Endpoint != "" {
					endpoints = append(endpoints, c.Workspace.Endpoint)
				}

				table.Append([]string{c.ID, c.Name, c.Project, c.Description})
				for _, s := range c.Stages {
					stage := CodeStreamPipelineStage{}
					mapstructure.Decode(s, &stage)
					// Loop through the Stage Tasks
					for n, t := range stage.Tasks {
						taskString := fmt.Sprintf("%v", t)
						rxp := regexp.MustCompile(`\$\{var\.(.*?)\}`) // Match ${var.name}
						variableMatches := rxp.FindAllStringSubmatch(taskString, -1)
						for _, v := range variableMatches {
							variables = append(variables, v[1])
						}
						task := CodeStreamPipelineTask{}
						mapstructure.Decode(t, &task)
						if len(task.Endpoints) > 0 {
							for _, e := range task.Endpoints {
								endpoints = append(endpoints, e)
							}
						}
						//PrettyPrint(task)
						if task.Type == "Pipeline" {
							pipelines = append(pipelines, task.Input.Pipeline)
						}
						if task.Type == "Custom" {
							customintegrations = append(customintegrations, task.Input.Name)
						}
						log.Infoln("-- [Task]", n, "(", task.Type, ")")
					}
				}
				if dependencies {
					variables = removeDuplicateStrings(variables)
					sort.Strings(variables)
					if len(variables) > 0 {
						log.Infoln(c.Name, "depends on Variables:", strings.Join(variables, ", "))
						for _, v := range variables {
							getVariable("", v, c.Project, exportPath)
						}
					}
					pipelines = removeDuplicateStrings(pipelines)
					sort.Strings(pipelines)
					if len(pipelines) > 0 {
						log.Infoln(c.Name, "depends on Pipelines:", strings.Join(pipelines, ", "))
						for _, p := range pipelines {
							getPipelines("", p, c.Project, filepath.Join(exportPath, "pipelines"))
						}
					}
					endpoints = removeDuplicateStrings(endpoints)
					sort.Strings(endpoints)
					if len(endpoints) > 0 {
						log.Infoln(c.Name, "depends on Endpoints:", strings.Join(endpoints, ", "))
						for _, e := range endpoints {
							getEndpoint("", e, c.Project, "", filepath.Join(exportPath, "endpoints"))
						}
					}
					customintegrations = removeDuplicateStrings(customintegrations)
					sort.Strings(customintegrations)
					if len(customintegrations) > 0 {
						log.Infoln(c.Name, "depends on Custom Integrations:", strings.Join(customintegrations, ", "))
						for _, ci := range customintegrations {
							getCustomIntegration("", ci)
						}
					}
				}
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
# Enable/Disable/Release:
cs-cli update pipeline --id d0185f04-2e87-4f3c-b6d7-ee58abba3e92 --state enabled/disabled/released
# Update from YAML
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
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if state != "" {
			response, err := patchPipeline(id, `{"state":"`+state+`"}`)
			if err != nil {
				log.Errorln("Unable to update Code Stream Pipeline: ", err)
			}
			log.Infoln("Setting pipeline " + response.Name + " to " + state)
		}

		yamlFilePaths := getYamlFilePaths(importPath)
		if len(yamlFilePaths) == 0 {
			log.Warnln("No YAML files were found in", importPath)
		}
		for _, yamlFilePath := range yamlFilePaths {
			yamlFileName := filepath.Base(yamlFilePath)
			err := importYaml(yamlFilePath, "apply", "", "endpoint")
			if err != nil {
				log.Warnln("Failed to import", yamlFilePath, "as Pipeline", err)
			}
			fmt.Println("Imported", yamlFileName, "successfully - Pipeline updated.")
		}
	},
}

// createPipelineCmd represents the pipeline create command
var createPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Create a Pipeline",
	Long: `Create a Pipeline by importing a YAML specification.
	
# Create from YAML
cs-cli create pipeline --importPath "/Users/sammcgeown/Desktop/pipelines/SSH Exports.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}
		yamlFilePaths := getYamlFilePaths(importPath)
		if len(yamlFilePaths) == 0 {
			log.Warnln("No YAML files were found in", importPath)
		}
		for _, yamlFilePath := range yamlFilePaths {
			yamlFileName := filepath.Base(yamlFilePath)
			err := importYaml(yamlFilePath, "create", project, "pipeline")
			if err != nil {
				log.Warnln("Failed to import", yamlFilePath, "as Pipeline", err)
			} else {
				fmt.Println("Imported", yamlFileName, "successfully - Pipeline created.")
			}
		}
	},
}

// deletePipelineCmd represents the delete pipeline command
var deletePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Delete a Pipeline",
	Long: `Delete a Pipeline with a specific ID, Name or by Project

# Delete by ID
cs-cli delete pipeline --id "pipeline ID"

# Delete by Name
cs-cli delete pipeline --name "My Pipeline"

# Delete by Name and Project
cs-cli delete pipeline --name "My Pipeline" --project "My Project"

# Delete all pipelines in Project
cs-cli delete pipeline --project "My Project"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}
		if id != "" {
			response, err := deletePipeline(id)
			if err != nil {
				log.Errorln("Delete Pipeline failed:", err)
			}
			log.Infoln("Pipeline with id " + response.ID + " deleted")
		} else if project != "" {
			response, err := deletePipelineInProject(project)
			if err != nil {
				log.Errorln("Delete Pipelines in "+project+" failed:", err)
			} else {
				log.Infoln(len(response), "Pipelines deleted")
			}

		}

	},
}

func init() {
	// Get
	getCmd.AddCommand(getPipelineCmd)
	getPipelineCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to list executions for")
	getPipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	getPipelineCmd.Flags().StringVarP(&project, "project", "p", "", "List pipeline in project")
	getPipelineCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getPipelineCmd.Flags().BoolVarP(&printForm, "form", "f", false, "Return pipeline inputs form(s)")
	getPipelineCmd.Flags().BoolVarP(&printJson, "json", "", false, "Return JSON formatted Pipeline(s)")
	getPipelineCmd.Flags().BoolVarP(&dependencies, "exportDependencies", "", false, "Export Pipeline dependencies (Endpoint, Pipelines, Variables, Custom Integrations)")

	// Create
	createCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().StringVarP(&importPath, "importPath", "", "", "YAML configuration file to import")
	createPipelineCmd.Flags().StringVarP(&project, "project", "p", "", "Manually specify the Project in which to create the Pipeline (overrides YAML)")
	createPipelineCmd.MarkFlagRequired("importPath")
	// Update
	updateCmd.AddCommand(updatePipelineCmd)
	updatePipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	updatePipelineCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	updatePipelineCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the pipeline (ENABLED|DISABLED|RELEASED")
	// Delete
	deleteCmd.AddCommand(deletePipelineCmd)
	deletePipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Pipeline to delete")
	deletePipelineCmd.Flags().StringVarP(&project, "project", "p", "", "Delete all Pipelines in the specified Project")

}
