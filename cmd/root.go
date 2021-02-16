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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var apiKey string
var server string
var id string
var name string
var project string
var typename string
var value string
var description string
var status string
var username string
var password string
var domain string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cs-cli",
	Short: "CLI Interface for VMware vRealize Automation Code Stream",
	Long:  `Command line interface for VMware vRealize Automation Code Stream`,
}

// Execute is the main process
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cs-cli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cs-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cs-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	var currentEndpointName = viper.GetString("currentEndpointName")
	apiKey = viper.GetString("endpoint." + currentEndpointName + ".apiKey")
	server = viper.GetString("endpoint." + currentEndpointName + ".server")

	// If the apiKey is not set or testAccesToken returns false
	if apiKey == "" || testAccessToken() == false {
		// Authenticate
		accessToken, authError := authenticate(viper.GetString("endpoint."+currentEndpointName+".server"), viper.GetString("endpoint."+currentEndpointName+".username"), viper.GetString("endpoint."+currentEndpointName+".password"), viper.GetString("endpoint."+currentEndpointName+".domain"))
		if authError != nil {
			fmt.Println("Authentication failed", authError.Error())
			os.Exit(1)
		}
		viper.Set("endpoint."+currentEndpointName+".apiKey", accessToken)
		viper.WriteConfig()
		apiKey = viper.GetString("endpoint." + currentEndpointName + ".apiKey")
	}
}
