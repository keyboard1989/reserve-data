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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/blockchain/nonce"
	"github.com/KyberNetwork/reserve-data/cmd/configuration"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/http"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
)

var noAuthEnable bool
var servPort int = 8000

func loadTimestamp(path string) []uint64 {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	timestamp := []uint64{}
	err = json.Unmarshal(raw, &timestamp)
	if err != nil {
		panic(err)
	}
	return timestamp
}

func serverStart(cmd *cobra.Command, args []string) {
	log.Println(args)
	log.Println(noAuthEnable)
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	var config *configuration.Config
	//get configuration
	switch kyberENV {
	case "mainnet", "production":
		log.Printf("Running in production mode")
		config = configuration.GetConfigForMainnet()
		break
	case "staging":
		log.Printf("Running in staging mode")
		config = configuration.GetConfigForStaging()
		break
	case "simulation":
		log.Printf("Running in simulation mode")
		config = configuration.GetConfigForSimulation()
		break
	case "kovan":
		log.Printf("Running in kovan mode")
		config = configuration.GetConfigForKovan()
		break
	case "ropsten":
		log.Printf("Running in ropsten mode")
		config = configuration.GetConfigForRopsten()
		break
	case "dev":
		log.Printf("Running in dev mode")
		config = configuration.GetConfig(configuration.ConfigPaths["dev"], configuration.NewDevExchangePool, !noAuthEnable)
	default:
		log.Printf("Running in dev mode")
		config = configuration.GetConfig(configuration.ConfigPaths["dev"], configuration.NewDevExchangePool, !noAuthEnable)
	}

	//set log file
	logPath := "/go/src/github.com/KyberNetwork/reserve-data/cmd/log.log"
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Couldn't open log file: %v", err)
	}
	// write to mutiple location :stdout and log path
	mw := io.MultiWriter(os.Stdout, f)
	defer f.Close()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(mw)

	//get fetcher based on config and ENV == stimulation.
	fetcher := fetcher.NewFetcher(
		config.FetcherStorage,
		config.FetcherRunner,
		config.ReserveAddress,
		os.Getenv("KYBER_ENV") == "simulation",
	)

	//set static field supportExchange from common...
	for _, ex := range config.Exchanges {
		common.SupportedExchanges[ex.ID()] = ex
	}
	for _, ex := range config.FetcherExchanges {
		fetcher.AddExchange(ex)
	}

	//set client & endpoint
	client, err := rpc.Dial(config.EthereumEndpoint)
	if err != nil {
		panic(err)
	}
	infura := ethclient.NewClient(client)
	bkclients := map[string]*ethclient.Client{}
	for _, ep := range config.BackupEthereumEndpoints {
		bkclient, err := ethclient.Dial(ep)
		if err != nil {
			log.Printf("Cannot connect to %s, err %s. Ignore it.", ep, err)
		} else {
			bkclients[ep] = bkclient
		}
	}

	// nonceCorpus := nonce.NewAutoIncreasing(infura, fileSigner)
	nonceCorpus := nonce.NewTimeWindow(infura, config.BlockchainSigner)
	nonceDeposit := nonce.NewTimeWindow(infura, config.DepositSigner)
	//set block chain
	bc, err := blockchain.NewBlockchain(
		client,
		infura,
		bkclients,
		config.WrapperAddress,
		config.PricingAddress,
		config.FeeBurnerAddress,
		config.NetworkAddress,
		config.ReserveAddress,
		config.BlockchainSigner,
		config.DepositSigner,
		nonceCorpus,
		nonceDeposit,
	)
	if err != nil {
		panic(err)
	}

	for _, token := range config.SupportedTokens {
		bc.AddToken(token)
	}
	err = bc.LoadAndSetTokenIndices()
	if err != nil {
		fmt.Printf("Can't load and set token indices: %s\n", err)
	} else {
		fetcher.SetBlockchain(bc)
		app := data.NewReserveData(
			config.DataStorage,
			fetcher,
		)
		app.Run()
		core := core.NewReserveCore(bc, config.ActivityStorage, config.ReserveAddress)
		server := http.NewHTTPServer(
			app, core,
			config.MetricStorage,
			":8000",
			config.EnableAuthentication,
			config.AuthEngine,
		)

		server.Run()

	}
}

// This represents the base command when called without any subcommands
var startServer = &cobra.Command{
	Use:   "serve ",
	Short: "initiate the server with specific config",
	Long: `Start reserve-data core server with preset Environment and
Allow overwriting some parameter`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: serverStart,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	startServer.Flags().BoolVarP(&noAuthEnable, "noauth", "", false, "enable authentication")
	startServer.Flags().IntVarP(&servPort, "port", "p", 8000, "server port")

}
