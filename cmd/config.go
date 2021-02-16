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

// getEndpointCmd represents the get-endpoint command
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
			if endpoint == nil {
				fmt.Println("Endpoint not found.")
			} else {
				PrettyPrint(endpoint)
			}
		} else {
			var endpoints = viper.GetStringMapString("endpoint")
			for key := range endpoints {
				fmt.Println(key)
			}
		}
	},
}

// setEndpointCmd represents the set-endpoint command
var setEndpointCmd = &cobra.Command{
	Use:   "set-endpoint",
	Short: "Creates a new endpoint config",
	Long: `Creates a new endpoint config

Examples:
	cs-cli config set-endpoint --name vra-test-ga --server vra8-test-ga.cmbu.local --username test-user --password VMware1! --domain cmbu.local
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	// current-endpoint
	configCmd.AddCommand(currentEndpointCmd)
	// get-endpoint
	configCmd.AddCommand(getEndpointCmd)
	getEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Display the endpoint with this name")
	// set-endpoint
	configCmd.AddCommand(setEndpointCmd)
	setEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Display the endpoint with this name")
	setEndpointCmd.Flags().StringVarP(&server, "server", "s", "", "Display the endpoint with this name")
	setEndpointCmd.Flags().StringVarP(&username, "username", "u", "", "Display the endpoint with this name")
	setEndpointCmd.Flags().StringVarP(&password, "password", "p", "", "Display the endpoint with this name")
	setEndpointCmd.Flags().StringVarP(&domain, "domain", "d", "", "Display the endpoint with this name")
	setEndpointCmd.MarkFlagRequired("name")
	setEndpointCmd.MarkFlagRequired("server")
	setEndpointCmd.MarkFlagRequired("username")
	setEndpointCmd.MarkFlagRequired("password")

}
