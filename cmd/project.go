/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getProjectCommand represents the project command
var getProjectCommand = &cobra.Command{
	Use:   "project",
	Short: "Get Projects",
	Long:  `Get Code Stream Projects`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		response, err := getProject(id, name)
		if err != nil {
			log.Errorln("Unable to get Code Stream Projects: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Warnln("No results found")
		}

		// Print result table
		table := tablewriter.NewWriter(os.Stdout)
		// pipelineTable := tablewriter.NewWriter(os.Stdout)
		// variableTable := tablewriter.NewWriter(os.Stdout)
		// endpointTable := tablewriter.NewWriter(os.Stdout)
		// table.SetHeader([]string{"Id", "Name", "Description"})
		for _, p := range response {
			table.Append([]string{p.ID, p.Name, p.Description})
			if exportPath != "" {
				tmpDir, err := ioutil.TempDir(os.TempDir(), "cs-cli-*")
				if err != nil {
					log.Fatalln(err)
				}
				zipFile := filepath.Join(exportPath, p.Name+".zip")
				var zipFiles []string
				log.Debugln(zipFile)
				pipelines, _ := getPipelines("", "", p.Name, filepath.Join(tmpDir, p.Name, "pipelines"))
				//pipelineTable.SetHeader([]string{"Id", "Name", "Project", "Description"})
				for _, c := range pipelines {
					zipFiles = append(zipFiles, filepath.Join(tmpDir, p.Name, "pipelines", c.Name+".yaml"))
					//pipelineTable.Append([]string{c.ID, c.Name, c.Project, c.Description})
				}
				variables, _ := getVariable("", "", p.Name, filepath.Join(tmpDir, p.Name))
				//variableTable.SetHeader([]string{"Id", "Name", "Project", "Description"})
				if len(variables) > 0 {
					zipFiles = append(zipFiles, filepath.Join(tmpDir, p.Name, "variables.yaml"))
				}
				// for _, c := range variables {
				// 	//variableTable.Append([]string{c.ID, c.Name, c.Project, c.Description})
				// }
				endpoints, _ := getEndpoint("", "", p.Name, "", filepath.Join(tmpDir, p.Name, "endpoints"))
				//endpointTable.SetHeader([]string{"ID", "Name", "Project", "Type", "Description"})
				for _, c := range endpoints {
					zipFiles = append(zipFiles, filepath.Join(tmpDir, p.Name, c.Name+".yaml"))
					//endpointTable.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description})
				}
				if err := ZipFiles(zipFile, zipFiles, tmpDir); err != nil {
					log.Fatalln(err)
				}
			}
		}
		// fmt.Println("Project")
		table.Render()
		// fmt.Println("Pipelines")
		// pipelineTable.Render()
		// fmt.Println("Variables")
		// variableTable.Render()
		// fmt.Println("Endpoints")
		// endpointTable.Render()
	},
}

func init() {
	// Get
	getCmd.AddCommand(getProjectCommand)
	getProjectCommand.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to list executions for")
	getProjectCommand.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	getProjectCommand.Flags().StringVarP(&exportPath, "exportpath", "", "", "Path to export projects and contents")

}
