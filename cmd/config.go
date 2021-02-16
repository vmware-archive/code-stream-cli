/*
Package cmd Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// currentEndpointCmd represents the current-endpoint command
var currentEndpointCmd = &cobra.Command{
	Use:   "current-endpoint",
	Short: "Display the current-endpoint",
	Long: `Displays the current-endpoint

Examples:
	# Display the current-endpoint
	cs-cli config current-endpoint
`,
	Run: func(cmd *cobra.Command, args []string) {
		var currentEndpointName = viper.GetString("currentEndpointName")
		if currentEndpointName != "" {
			fmt.Println(currentEndpointName)
		}
	},
}

// getEndpointCmd represents the current-endpoint command
var getEndpointCmd = &cobra.Command{
	Use:   "get-endpoint",
	Short: "Display available endpoint configs",
	Long: `Displays a list of the available endpoint configs

Examples:
	cs-cli config get-endpoint
`,
	Run: func(cmd *cobra.Command, args []string) {
		if name != "" {
			var endpoint = viper.Get("endpoint." + name)
			PrettyPrint(endpoint)
		} else {
			var endpoints = viper.GetStringMapString("endpoint")
			for key := range endpoints {
				fmt.Println(key)
			}
		}
	},
}

func init() {
	// Get Variable
	configCmd.AddCommand(currentEndpointCmd)
	configCmd.AddCommand(getEndpointCmd)
	getEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Display the endpoint with this name")

}
