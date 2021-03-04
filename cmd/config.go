/*
Package cmd Copyright © 2021 Sam McGeown <smcgeown@vmware.com>

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
	"log"

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
			log.Println(currentTargetName)
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
			log.Println("Target not found! Current target is", viper.GetString("currentTargetName"))
			return
		}
		viper.Set("currentTargetName", name)
		viper.WriteConfig()
		log.Println("Current target: ", name)
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
				log.Println("Target not found.")
			} else {
				PrettyPrint(target)
			}
		} else {
			var targets = viper.GetStringMapString("target")
			for key := range targets {
				log.Println(key)
			}
		}
	},
}

var (
	newTargetName string
	newServer     string
	newUsername   string
	newPassword   string
	newDomain     string
	newAPIToken   string
)

// setTargetCmd represents the set-target command
var setTargetCmd = &cobra.Command{
	Use:   "set-target",
	Short: "Creates or updates a target config",
	Long: `Creates or updates a target configuration.

Examples:
	cs-cli config set-target --name vra-test-ga --server vra8-test-ga.cmbu.local --username test-user --password VMware1! --domain cmbu.local
	cs-cli config set-target --name vrac-org --server api.mgmt.cloud.vmware.com --apitoken JhbGciOiJSUzI1NiIsImtpZCI6IjEzNjY3NDcwMTA2Mzk2MTUxNDk0In0
`, Args: func(cmd *cobra.Command, args []string) error {
		// if apiToken != "" && server != "" && username == "" && password == "" {
		// 	return nil
		// } else if apiToken == "" && server != "" && username != "" && password != "" {
		// 	return nil
		// }
		// return errors.New("Incorrect combination of flags, please use  --server and --apitoken for vRealize Automation Cloud, or --server, --username, --password and --domain (optional) for vRealize Automation 8.x")
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("target." + newTargetName) {
			log.Println("Updating", newTargetName)
		} else {
			log.Println("Creating new target", newTargetName)
		}
		log.Println("Use `cs-cli config use-target --name " + newTargetName + "` to use this target")
		if newServer != "" {
			viper.Set("target."+newTargetName+".server", newServer)
		}
		if newUsername != "" {
			viper.Set("target."+newTargetName+".username", newUsername)
		}
		if newPassword != "" {
			viper.Set("target."+newTargetName+".password", newPassword)
		}
		if newDomain != "" {
			viper.Set("target."+newTargetName+".domain", newDomain)
		}
		if newAPIToken != "" {
			viper.Set("target."+newTargetName+".apitoken", newAPIToken)
		}
		err := viper.WriteConfig()
		if err != nil {
			log.Println(err)
		}
		var target = viper.Get("target." + newTargetName)
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
// 				log.Println(err)
// 			}
// 			log.Println("Target deleted.")
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
	setTargetCmd.Flags().StringVarP(&newTargetName, "name", "n", "", "Name of the target configuration")
	setTargetCmd.Flags().StringVarP(&newServer, "server", "s", "", "Server FQDN of the vRealize Automation instance")
	setTargetCmd.Flags().StringVarP(&newUsername, "username", "u", "", "Username to authenticate")
	setTargetCmd.Flags().StringVarP(&newPassword, "password", "p", "", "Password to authenticate")
	setTargetCmd.Flags().StringVarP(&newDomain, "domain", "d", "", "Domain to authenticate (not required for System Domain)")
	setTargetCmd.Flags().StringVarP(&newAPIToken, "apitoken", "a", "", "API token for vRealize Automation Cloud")
	setTargetCmd.MarkFlagRequired("name")
	// delete-target
	// configCmd.AddCommand(deleteTargetCmd)
	// deleteTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the target configuration")
	// deleteTargetCmd.MarkFlagRequired("name")
}
