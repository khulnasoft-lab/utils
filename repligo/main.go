package main

import (
	"fmt"
	"github.com/khulnasoft-lab/utils/repligo"
	"github.com/spf13/viper"
	"os"
)

var version = ""

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" && version != "" {
		fmt.Println(version)
		os.Exit(0)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("repligo")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("/etc/repligo/")
	viper.AddConfigPath("$HOME/.repligo")

	// AWS_PROFILE or REPLIGO_PROFILE can set the profile
	viper.BindEnv("profile", "AWS_PROFILE");

	for _, ext := range []string{"json", "yaml", "yml"} {
		viper.SetConfigType(ext)
		err := viper.ReadInConfig() // Find and read config files
		if err != nil { // Handle errors reading the config file
			if _, ok := err.(viper.ConfigParseError); ok {
				panic(fmt.Errorf("Fatal error reading " + ext + " config file: %s \n", err))
			}
		}
	}

	for _, extrafile := range []string{"$HOME/.aws/config", "$HOME/.aws/credentials"} {

		viper.SetConfigType("yaml") // TODO: this should be INI
		viper.SetConfigFile(extrafile) // name of config file (without extension)
		err := viper.ReadInConfig() // Find and read the config file

		if err != nil { // Handle errors reading the config file
			if _, ok := err.(viper.ConfigParseError); ok {
				panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
			}
		}
	}

	//TODO: Investigate: do we set the max procs somewhere?
	cmd.Execute()
}
