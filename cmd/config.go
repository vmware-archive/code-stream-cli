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

// currentTargetCmd represents the current-target command
var currentTargetCmd = &cobra.Command{
	Use:   "current-target",
	Short: "Display the current-target",
	Long: `Displays the current-target

Examples:
	# Display the current-target
	cs-cli config current-target
`,
	Run: func(cmd *cobra.Command, args []string) {
		var currentTargetName = viper.GetString("currentTargetName")
		if currentTargetName != "" {
			fmt.Println(currentTargetName)
		}
	},
}

// useTargetCmd represents the use-target command
var useTargetCmd = &cobra.Command{
	Use:   "use-target",
	Short: "Set the current target",
	Long: `Set the current target

Examples:
	# Display the current-target
	cs-cli config use-target --name vra8-test-ga
`,
	Run: func(cmd *cobra.Command, args []string) {
		var target = viper.Get("target." + name)
		if target == nil {
			fmt.Println("Target not found! Current target is", viper.GetString("currentTargetName"))
			return
		}
		viper.Set("currentTargetName", name)
		viper.WriteConfig()
		fmt.Println("Current target: ", name)
	},
}

// getConfigTargetCmd represents the get-target command
var getConfigTargetCmd = &cobra.Command{
	Use:   "get-target",
	Short: "Display available target configs",
	Long: `Displays a list of the available target configs

Examples:
	cs-cli config get-target
`,
	Run: func(cmd *cobra.Command, args []string) {
		if name != "" {
			var target = viper.Get("target." + name)
			if target == nil {
				fmt.Println("Target not found.")
			} else {
				PrettyPrint(target)
			}
		} else {
			var targets = viper.GetStringMapString("target")
			for key := range targets {
				fmt.Println(key)
			}
		}
	},
}

var newServer string

// setTargetCmd represents the set-target command
var setTargetCmd = &cobra.Command{
	Use:   "set-target",
	Short: "Creates or updates an target config",
	Long: `Creates or updates an target config

Examples:
	cs-cli config set-target --name vra-test-ga --server vra8-test-ga.cmbu.local --username test-user --password VMware1! --domain cmbu.local
`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("target." + name) {
			fmt.Println("Updating", name)
		} else {
			fmt.Println("Creating new target", name)
		}
		fmt.Println("Use `cs-cli config use-target --name " + name + "` to use this target")
		if newServer != "" {
			viper.Set("target."+name+".server", newServer)
		}
		if username != "" {
			viper.Set("target."+name+".username", username)
		}
		if password != "" {
			viper.Set("target."+name+".password", password)
		}
		if domain != "" {
			viper.Set("target."+name+".domain", domain)
		}
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
		var target = viper.Get("target." + name)
		PrettyPrint(target)
	},
}

// deleteTargetCmd represents the set-target command
// var deleteTargetCmd = &cobra.Command{
// 	Use:   "delete-target",
// 	Short: "Deletes an target config",
// 	Long: `Deletes an target config

// Examples:
// 	cs-cli config delete-target --name vra-test-ga
// `,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if viper.IsSet("target." + name) {
// 			err := Unset(name)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			fmt.Println("Target deleted.")
// 		}
// 	},
// }

func init() {
	// current-target
	configCmd.AddCommand(currentTargetCmd)
	// use-target
	configCmd.AddCommand(useTargetCmd)
	useTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Use the target with this name")
	useTargetCmd.MarkFlagRequired("name")
	// get-target
	configCmd.AddCommand(getConfigTargetCmd)
	getConfigTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Display the target with this name")
	// set-target
	configCmd.AddCommand(setTargetCmd)
	setTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the target configuration")
	setTargetCmd.Flags().StringVarP(&newServer, "server", "s", "", "Server FQDN of the target")
	setTargetCmd.Flags().StringVarP(&username, "username", "u", "", "Username to authenticate with the target")
	setTargetCmd.Flags().StringVarP(&password, "password", "p", "", "Password to authenticate with the target")
	setTargetCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain to authenticate with the target (not required for System Domain)")
	setTargetCmd.MarkFlagRequired("name")
	// delete-target
	// configCmd.AddCommand(deleteTargetCmd)
	// deleteTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the target configuration")
	// deleteTargetCmd.MarkFlagRequired("name")
}
