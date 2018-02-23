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
	"log"

	"github.com/spf13/cobra"
)

var baseURL string

func compareratestart(cmd *cobra.Command, args []string) {
	log.Println(baseURL)
}

// This represents the base command when called without any subcommands
var compareRates = &cobra.Command{
	Use:   "compare ",
	Short: "compare rate from get_all_rates to setRate activities",
	Long: `Start reserve-data core server with preset Environment and
Allow overwriting some parameter`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: compareratestart,
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	compareRates.Flags().StringVar(&baseURL, "url", "https://ropsten-core.kyber.network", "base URL, default to https://ropsten-core.kyber.network")
	compareRates.MarkFlagRequired("url")
	RootCmd.AddCommand(compareRates)
}
