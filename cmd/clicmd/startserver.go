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
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var noAuthEnable bool
var servPort int = 8000
var addressOW [5]string
var endpointOW string

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

// GetConfigFromENV: From ENV variable and overwriting instruction, build the config
func GetConfigFromENV(kyberENV string, addressOW [5]string) *configuration.Config {
	var config *configuration.Config
	config = configuration.GetConfig(kyberENV,
		!noAuthEnable,
		addressOW,
		endpointOW)
	return config
}

//set config log
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

	c := cron.New()
	c.AddFunc("@daily", func() { logger.Rotate() })
	c.Start()
}

func serverStart(cmd *cobra.Command, args []string) {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	configLog()
	//get configuration from ENV variable
	kyberENV := os.Getenv("KYBER_ENV")
	if kyberENV == "" {
		kyberENV = "dev"
	}
	config := GetConfigFromENV(kyberENV, addressOW)

	//get fetcher based on config and ENV == stimulation.
	fetcher := fetcher.NewFetcher(
		config.FetcherStorage,
		config.FetcherRunner,
		config.ReserveAddress,
		kyberENV == "simulation",
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

	//nonceCorpus := nonce.NewAutoIncreasing(infura, fileSigner)
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
		servPortStr := fmt.Sprintf(":%d", servPort)
		server := http.NewHTTPServer(
			app, core,
			config.MetricStorage,
			servPortStr,
			config.EnableAuthentication,
			config.AuthEngine,
			kyberENV,
		)

		server.Run()

	}
}

var startServer = &cobra.Command{
	Use:   "server ",
	Short: "initiate the server with specific config",
	Long: `Start reserve-data core server with preset Environment and
Allow overwriting some parameter`,
	Run: serverStart,
}

func init() {
	// start server flags.
	startServer.Flags().BoolVarP(&noAuthEnable, "noauth", "", false, "disable authentication")
	startServer.Flags().IntVarP(&servPort, "port", "p", 8000, "server port")
	startServer.Flags().StringVar(&addressOW[0], "wrapperAddr", "", "wrapper Address, default to configuration file")
	startServer.Flags().StringVar(&addressOW[1], "reserveAddr", "", "reserve Address, default to configuration file")
	startServer.Flags().StringVar(&addressOW[2], "pricingAddr", "", "pricing Address, default to configuration file")
	startServer.Flags().StringVar(&addressOW[3], "burnerAddr", "", "burner Address, default to configuration file")
	startServer.Flags().StringVar(&addressOW[4], "networkAddr", "", "network Address, default to configuration file")
	startServer.Flags().StringVar(&endpointOW, "endpoint", "", "endpoint, default to configuration file")
	RootCmd.AddCommand(startServer)
}
