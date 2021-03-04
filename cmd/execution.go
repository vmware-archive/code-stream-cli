package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var nested bool
var inputs string
var comments string
var inputPath string

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
		ensureTargetConnection()

		response, err := getExecutions(id, status, name, nested)
		if err != nil {
			log.Println("Unable to get executions: ", err)
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
			table.SetHeader([]string{"Id", "Name", "Project", "Status", "Message"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name + "#" + fmt.Sprint(c.Index), c.Project, c.Status, c.StatusMessage})
			}
			table.Render()
		}

	},
}

// delExecutionCmd represents the executions command
var delExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Delete an execution",
	Long: `Delete an execution with a specific execution ID
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()

		response, err := deleteExecution(id)
		if err != nil {
			log.Println("Unable to delete execution: ", err)
		}
		log.Println("Execution with id " + response.ID + " deleted")

	},
}

// createExecutionCmd represents the executions command
var createExecutionCmd = &cobra.Command{
	Use:   "execution",
	Short: "Create an execution",
	Long: `Create an execution with a specific pipeline ID and form payload.
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()

		response, err := createExecution(id, inputs, comments)
		if err != nil {
			log.Println("Unable to create execution: ", err)
		}
		log.Println("Execution " + response.ExecutionLink + " created")

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
	// Create
	createCmd.AddCommand(createExecutionCmd)
	createExecutionCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to execute")
	createExecutionCmd.Flags().StringVarP(&inputs, "inputs", "", "", "JSON form inputs")
	createExecutionCmd.Flags().StringVarP(&inputPath, "inputPath", "", "", "JSON input file")
	createExecutionCmd.Flags().StringVarP(&comments, "comments", "", "", "Execution comments")
	createExecutionCmd.MarkFlagRequired("id")
}
