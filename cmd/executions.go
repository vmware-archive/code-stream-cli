/*
Package cmd ...
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var n string
var i string

// executionsCmd represents the executions command
var executionsCmd = &cobra.Command{
	Use:   "executions",
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
		for _, x := range response {
			fmt.Println(x.Name)
		}
	},
}

func init() {
	getCmd.AddCommand(executionsCmd)
	executionsCmd.Flags().StringVarP(&n, "name", "n", "", "Name of the pipeline to list executions for")
	executionsCmd.Flags().StringVarP(&i, "id", "i", "", "ID of the pipeline to list executions for")
}
