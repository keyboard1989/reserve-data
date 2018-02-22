package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

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
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

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

func configLog() {
	logger := &lumberjack.Logger{
		Filename: "/go/src/github.com/KyberNetwork/reserve-data/log/core.log",
		// MaxSize:  1, // megabytes
		MaxBackups: 0,
		MaxAge:     0, //days
		// Compress:   true, // disabled by default
	}

	mw := io.MultiWriter(os.Stdout, logger)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(mw)

	rotateTicker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			<-rotateTicker.C
			logger.Rotate()
		}
	}()
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	configLog()

	var config *configuration.Config
	env := os.Getenv("KYBER_ENV")
	switch env {
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
		config = configuration.GetConfigForDev()
	default:
		log.Printf("Running in dev mode")
		config = configuration.GetConfigForDev()
	}

	fetcher := fetcher.NewFetcher(
		config.FetcherStorage,
		config.FetcherRunner,
		config.ReserveAddress,
		os.Getenv("KYBER_ENV") == "simulation",
	)
	for _, ex := range config.Exchanges {
		common.SupportedExchanges[ex.ID()] = ex
	}
	for _, ex := range config.FetcherExchanges {
		fetcher.AddExchange(ex)
	}
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
			env,
		)
		server.Run()
	}
}
