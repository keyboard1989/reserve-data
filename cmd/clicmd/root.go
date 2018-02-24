// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addressConfigFile *string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "./cmd ",
	Short: "entry point to the application, required KYBER_ENV and KYBER_EXCHANGES as environment variables",

	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.AddCommand(startServer)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	addressConfigFile = RootCmd.PersistentFlags().String("addressconfig", "", "default config file name")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("verbose", "v", false, "verbose mode enable")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if (*addressConfigFile) != "" { // enable ability to specify config file via fl
		viper.SetConfigName(*addressConfigFile)
	}

	viper.SetConfigType("json")
	viper.AddConfigPath("/go/src/github.com/KyberNetwork/reserve-data/cmd/")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %v \n", viper.ConfigFileUsed())
	}
}
