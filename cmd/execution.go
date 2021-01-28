/*
Package cmd ...
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var n string
var i string

// executionCmd represents the executions command
var executionCmd = &cobra.Command{
	Use:   "execution",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := getExecutions(i)
		if err != nil {
			fmt.Print("Unable to get executions: ", err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "Name", "Project", "Status", "Message"})

		for _, c := range response {
			table.Append([]string{c.ID, c.Name + "#" + fmt.Sprint(c.Index), c.Project, c.Status, c.StatusMessage})
		}
		table.Render()

	},
}

func init() {
	getCmd.AddCommand(executionCmd)
	executionCmd.Flags().StringVarP(&n, "name", "n", "", "Name of the pipeline to list executions for")
	executionCmd.Flags().StringVarP(&i, "id", "i", "", "ID of the pipeline to list executions for")
}
