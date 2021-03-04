package cmd

import (
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mrz1836/go-sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Configuration
	cfgFile           string
	currentTargetName string
	targetConfig      config
	// Flags
	id                string
	name              string
	project           string
	typename          string
	value             string
	description       string
	status            string
	exportFile        string
	importFile        string
)

var qParams = map[string]string{
	"apiVersion": "2019-10-17",
}

type config struct {
	domain      string
	password    string
	server      string
	username    string
	apitoken    string
	accesstoken string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cs-cli",
	Short: "CLI Interface for VMware vRealize Automation Code Stream",
	Long:  `Command line interface for VMware vRealize Automation Code Stream`,
}

// Execute is the main process
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cs-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&targetConfig.server, "server", "", "vRealize Automation Server to target")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigName(".cs-cli")
	viper.SetConfigType("yaml")

	// Bind ENV variables
	viper.SetEnvPrefix("cs")
	viper.AutomaticEnv()

	// If we're using ENV variables
	if viper.Get("server") != nil {
		log.Println("Using ENV variables")
		targetConfig = config{
			server:      sanitize.URL(viper.GetString("server")),
			username:    viper.GetString("username"),
			password:    viper.GetString("password"),
			domain:      viper.GetString("domain"),
			apitoken:    viper.GetString("apitoken"),
			accesstoken: viper.GetString("accesstoken"),
		}
	} else {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("Fatal error config file:", err)
		}
		currentTargetName = viper.GetString("currentTargetName")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("Unable to load configuration:", err)
		}
		log.Println("Using config file:", viper.ConfigFileUsed())
		log.Println("Using config name:", currentTargetName)
		configuration := viper.Sub("target." + currentTargetName)
		if configuration == nil { // Sub returns nil if the key cannot be found
			log.Fatalln("Target configuration not found")
		}
		targetConfig = config{
			server:      sanitize.URL(configuration.GetString("server")),
			username:    configuration.GetString("username"),
			password:    configuration.GetString("password"),
			domain:      configuration.GetString("domain"),
			apitoken:    configuration.GetString("apitoken"),
			accesstoken: configuration.GetString("accesstoken"),
		}
	}
}
