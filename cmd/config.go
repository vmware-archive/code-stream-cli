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

// useEndpointCmd represents the use-endpoint command
var useEndpointCmd = &cobra.Command{
	Use:   "use-endpoint",
	Short: "Set the current endpoint",
	Long: `Set the current endpoint

Examples:
	# Display the current-endpoint
	cs-cli config use-endpoint --name vra8-test-ga
`,
	Run: func(cmd *cobra.Command, args []string) {
		var endpoint = viper.Get("endpoint." + name)
		if endpoint == nil {
			fmt.Println("Endpoint not found! Current endpoint is", viper.GetString("currentEndpointName"))
			return
		}
		viper.Set("currentEndpointName", name)
		viper.WriteConfig()
		fmt.Println("Current endpoint: ", name)
	},
}

// getConfigEndpointCmd represents the get-endpoint command
var getConfigEndpointCmd = &cobra.Command{
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

var newServer string

// setEndpointCmd represents the set-endpoint command
var setEndpointCmd = &cobra.Command{
	Use:   "set-endpoint",
	Short: "Creates or updates an endpoint config",
	Long: `Creates or updates an endpoint config

Examples:
	cs-cli config set-endpoint --name vra-test-ga --server vra8-test-ga.cmbu.local --username test-user --password VMware1! --domain cmbu.local
`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("endpoint." + name) {
			fmt.Println("Updating", name)
		} else {
			fmt.Println("Creating new endpoint", name)
		}
		fmt.Println("Use `cs-cli config use-endpoint --name " + name + "` to use this endpoint")
		if newServer != "" {
			viper.Set("endpoint."+name+".server", newServer)
		}
		if username != "" {
			viper.Set("endpoint."+name+".username", username)
		}
		if password != "" {
			viper.Set("endpoint."+name+".password", password)
		}
		if domain != "" {
			viper.Set("endpoint."+name+".domain", domain)
		}
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
		var endpoint = viper.Get("endpoint." + name)
		PrettyPrint(endpoint)
	},
}

// deleteEndpointCmd represents the set-endpoint command
// var deleteEndpointCmd = &cobra.Command{
// 	Use:   "delete-endpoint",
// 	Short: "Deletes an endpoint config",
// 	Long: `Deletes an endpoint config

// Examples:
// 	cs-cli config delete-endpoint --name vra-test-ga
// `,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if viper.IsSet("endpoint." + name) {
// 			err := Unset(name)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			fmt.Println("Endpoint deleted.")
// 		}
// 	},
// }

func init() {
	// current-endpoint
	configCmd.AddCommand(currentEndpointCmd)
	// use-endpoint
	configCmd.AddCommand(useEndpointCmd)
	useEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Use the endpoint with this name")
	useEndpointCmd.MarkFlagRequired("name")
	// get-endpoint
	configCmd.AddCommand(getConfigEndpointCmd)
	getConfigEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Display the endpoint with this name")
	// set-endpoint
	configCmd.AddCommand(setEndpointCmd)
	setEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the endpoint configuration")
	setEndpointCmd.Flags().StringVarP(&newServer, "server", "s", "", "Server FQDN of the endpoint")
	setEndpointCmd.Flags().StringVarP(&username, "username", "u", "", "Username to authenticate with the endpoint")
	setEndpointCmd.Flags().StringVarP(&password, "password", "p", "", "Password to authenticate with the endpoint")
	setEndpointCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain to authenticate with the endpoint (not required for System Domain)")
	setEndpointCmd.MarkFlagRequired("name")
	// delete-endpoint
	// configCmd.AddCommand(deleteEndpointCmd)
	// deleteEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the endpoint configuration")
	// deleteEndpointCmd.MarkFlagRequired("name")
}
