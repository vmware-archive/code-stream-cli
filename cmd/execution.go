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

var nested bool

// getExecutionCmd represents the executions command
var getExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Get executions of Code Stream Pipelines",
	Long: `List executions of Code Stream Pipelines by ID, Pipeline name, Project and Status
	Get only failed executions:
	  cs-cli get execution --status FAILED
	Get an execution by ID:
	  cs-cli get execution --id bb3f6aff-311a-45fe-8081-5845a529068d
	`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := getExecutions(id, status, name, nested)
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

// delExecutionCmd represents the executions command
var delExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		response, err := deleteExecution(id)
		if err != nil {
			fmt.Print("Unable to delete execution: ", err)
		}
		fmt.Println("Execution with id " + response.ID + " deleted")

	},
}

func init() {
	// Get
	getCmd.AddCommand(getExecutionCmd)
	getExecutionCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to list executions for")
	getExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the executions to list")
	getExecutionCmd.Flags().StringVarP(&status, "status", "s", "", "Filter executions by status (Completed|Waiting|Pausing|Paused|Resuming|Running)")
	getExecutionCmd.Flags().BoolVarP(&nested, "nested", "", false, "Include nested executions")
	// Delete
	deleteCmd.AddCommand(delExecutionCmd)
	delExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to delete")
	delExecutionCmd.MarkFlagRequired("id")
}
