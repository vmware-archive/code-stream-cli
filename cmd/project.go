/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

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

		response, err := getProject("", "", "")
		if err != nil {
			log.Println("Unable to get Code Stream Projects: ", err)
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
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Description"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Description})
			}
			table.Render()
		}
	},
}

func init() {
	// Get
	getCmd.AddCommand(getProjectCommand)

}
