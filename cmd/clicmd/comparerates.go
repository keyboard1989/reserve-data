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
	"os"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-data/cmd/comparerates"
	"github.com/spf13/cobra"
)

const (
	SleepTime time.Duration = 60 //sleep time for forever run mode
)

var baseURL string
var fromTime string
var toTime string

func compareratestart(cmd *cobra.Command, args []string) {
	kyberENV := os.Getenv("KYBER_ENV")
	if kyberENV == "" {
		kyberENV = "dev"
	}
	params := make(map[string]string)
	params["fromTime"] = fromTime
	params["toTime"] = toTime
	config := GetConfigFromENV(kyberENV, addressOW)
	if len(params["toTime"]) < 1 {
		log.Printf("There was no end time, go to foverer run mode...")
		for {
			params["toTime"] = strconv.FormatInt((time.Now().UnixNano() / int64(time.Millisecond)), 10)
			comparerates.DoQuery(baseURL, params, *config)
			time.Sleep(SleepTime * time.Second)
			params["fromTime"] = params["toTime"]
		}

	} else {
		log.Printf("Go to single query returning mode")
		comparerates.DoQuery(baseURL, params, *config)
	}

}

// This represents the base command when called without any subcommands
var compareRates = &cobra.Command{
	Use:   "compare ",
	Short: "compare rate from get_all_rates to setRate activities",
	Long:  `call get_all_rates and get_activities API to server and compare the rates between set_rate activites and get_rate activites from the same block, alert if the rate differ is >0.1%`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: compareratestart,
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	compareRates.Flags().StringVar(&baseURL, "url", "https://internal-mainnet-core.kyber.network", "base URL for API query")
	compareRates.Flags().StringVar(&fromTime, "from_time", "", "begining time for query, required params")
	compareRates.MarkFlagRequired("from_time")
	compareRates.Flags().StringVar(&toTime, "to_time", "", "end time of querying, if not set then the program will run until force quit")
	RootCmd.AddCommand(compareRates)
}
