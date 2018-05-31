package configuration

import (
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/settings"

	"github.com/KyberNetwork/reserve-data/cmd/configuration/mode"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/archive"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/datapruner"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/fetcher/http_runner"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/exchange/bittrex"
	"github.com/KyberNetwork/reserve-data/exchange/huobi"
	"github.com/KyberNetwork/reserve-data/http"
	"github.com/KyberNetwork/reserve-data/metric"
	"github.com/KyberNetwork/reserve-data/stat"
	"github.com/KyberNetwork/reserve-data/stat/statpruner"
	statstorage "github.com/KyberNetwork/reserve-data/stat/storage"
	"github.com/KyberNetwork/reserve-data/world"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type SettingPaths struct {
	settingPath           string
	feePath               string
	dataStoragePath       string
	analyticStoragePath   string
	statStoragePath       string
	logStoragePath        string
	rateStoragePath       string
	userStoragePath       string
	feeSetRateStoragePath string
	secretPath            string
	endPoint              string
	bkendpoints           []string
}

type Config struct {
	ActivityStorage      core.ActivityStorage
	DataStorage          data.Storage
	DataGlobalStorage    data.GlobalStorage
	StatStorage          stat.StatStorage
	AnalyticStorage      stat.AnalyticStorage
	UserStorage          stat.UserStorage
	LogStorage           stat.LogStorage
	RateStorage          stat.RateStorage
	FeeSetRateStorage    stat.FeeSetRateStorage
	FetcherStorage       fetcher.Storage
	FetcherGlobalStorage fetcher.GlobalStorage
	MetricStorage        metric.MetricStorage
	Archive              archive.Archive

	World                *world.TheWorld
	FetcherRunner        fetcher.FetcherRunner
	DataControllerRunner datapruner.StorageControllerRunner
	StatFetcherRunner    stat.FetcherRunner
	StatControllerRunner statpruner.ControllerRunner
	FetcherExchanges     []fetcher.Exchange
	Exchanges            []common.Exchange
	BlockchainSigner     blockchain.Signer
	DepositSigner        blockchain.Signer
	//IntermediatorSigner blockchain.Signer

	EnableAuthentication bool
	AuthEngine           http.Authentication

	EthereumEndpoint        string
	BackupEthereumEndpoints []string
	Blockchain              *blockchain.BaseBlockchain

	WrapperAddress     ethereum.Address
	PricingAddress     ethereum.Address
	ReserveAddress     ethereum.Address
	FeeBurnerAddress   ethereum.Address
	NetworkAddress     ethereum.Address
	WhitelistAddress   ethereum.Address
	SetRateAddress     ethereum.Address
	ThirdPartyReserves []ethereum.Address

	// etherscan api key (optional)
	EtherscanApiKey string

	ChainType string
	Setting   *settings.Settings
}

// GetStatConfig: load config to run stat server only
func (self *Config) AddStatConfig(settingPath SettingPaths, addressConfig common.AddressConfig) {
	networkAddr := ethereum.HexToAddress(addressConfig.Network)
	burnerAddr := ethereum.HexToAddress(addressConfig.FeeBurner)
	whitelistAddr := ethereum.HexToAddress(addressConfig.Whitelist)

	thirdpartyReserves := []ethereum.Address{}
	for _, address := range addressConfig.ThirdPartyReserves {
		thirdpartyReserves = append(thirdpartyReserves, ethereum.HexToAddress(address))
	}

	analyticStorage, err := statstorage.NewBoltAnalyticStorage(settingPath.analyticStoragePath)
	if err != nil {
		panic(err)
	}

	statStorage, err := statstorage.NewBoltStatStorage(settingPath.statStoragePath)
	if err != nil {
		panic(err)
	}

	logStorage, err := statstorage.NewBoltLogStorage(settingPath.logStoragePath)
	if err != nil {
		panic(err)
	}

	rateStorage, err := statstorage.NewBoltRateStorage(settingPath.rateStoragePath)
	if err != nil {
		panic(err)
	}

	userStorage, err := statstorage.NewBoltUserStorage(settingPath.userStoragePath)
	if err != nil {
		panic(err)
	}

	feeSetRateStorage, err := statstorage.NewBoltFeeSetRateStorage(settingPath.feeSetRateStoragePath)
	if err != nil {
		panic(err)
	}

	var statFetcherRunner stat.FetcherRunner
	var statControllerRunner statpruner.ControllerRunner
	if mode.Get() == mode.SIMULATION_MODE {
		if statFetcherRunner, err = http_runner.NewHttpRunner(http_runner.WithHttpRunnerPort(8002)); err != nil {
			panic(err)
		}
	} else {
		statFetcherRunner = stat.NewTickerRunner(
			5*time.Second,  // block fetching interval
			7*time.Second,  // log fetching interval
			10*time.Second, // rate fetching interval
			2*time.Second,  // tradelog processing interval
			2*time.Second)  // catlog processing interval
		statControllerRunner = statpruner.NewControllerTickerRunner(24 * time.Hour)

	}

	apiKey := GetEtherscanAPIKey(settingPath.secretPath)

	self.StatStorage = statStorage
	self.AnalyticStorage = analyticStorage
	self.UserStorage = userStorage
	self.LogStorage = logStorage
	self.RateStorage = rateStorage
	self.StatControllerRunner = statControllerRunner
	self.FeeSetRateStorage = feeSetRateStorage
	self.StatFetcherRunner = statFetcherRunner
	self.ThirdPartyReserves = thirdpartyReserves
	self.FeeBurnerAddress = burnerAddr
	self.NetworkAddress = networkAddr
	self.WhitelistAddress = whitelistAddr
	self.EtherscanApiKey = apiKey
}

func (self *Config) AddCoreConfig(settingPath SettingPaths, addressConfig common.AddressConfig, kyberENV string) {
	networkAddr := ethereum.HexToAddress(addressConfig.Network)
	burnerAddr := ethereum.HexToAddress(addressConfig.FeeBurner)
	whitelistAddr := ethereum.HexToAddress(addressConfig.Whitelist)

	feeConfig, err := common.GetFeeFromFile(settingPath.feePath)
	if err != nil {
		log.Fatalf("Fees file %s cannot found at: %s", settingPath.feePath, err)
	}

	minDepositPath := "/go/src/github.com/KyberNetwork/reserve-data/cmd/min_deposit.json"
	minDeposit, err := common.GetMinDepositFromFile(minDepositPath)
	if err != nil {
		log.Fatalf("Fees file %s cannot found at: %s", minDepositPath, err.Error())
	}

	dataStorage, err := storage.NewBoltStorage(settingPath.dataStoragePath)
	if err != nil {
		panic(err)
	}

	var fetcherRunner fetcher.FetcherRunner
	var dataControllerRunner datapruner.StorageControllerRunner
	if mode.Get() == mode.SIMULATION_MODE {
		if fetcherRunner, err = http_runner.NewHttpRunner(http_runner.WithHttpRunnerPort(8001)); err != nil {
			log.Fatalf("failed to create HTTP runner: %s", err.Error())
		}
	} else {
		fetcherRunner = fetcher.NewTickerRunner(
			7*time.Second,  // orderbook fetching interval
			5*time.Second,  // authdata fetching interval
			3*time.Second,  // rate fetching interval
			5*time.Second,  // block fetching interval
			10*time.Minute, // tradeHistory fetching interval
			10*time.Second, // global data fetching interval
		)
		dataControllerRunner = datapruner.NewStorageControllerTickerRunner(24 * time.Hour)
	}

	pricingSigner := PricingSignerFromConfigFile(settingPath.secretPath)
	depositSigner := DepositSignerFromConfigFile(settingPath.secretPath)

	self.ActivityStorage = dataStorage
	self.DataStorage = dataStorage
	self.DataGlobalStorage = dataStorage
	self.FetcherStorage = dataStorage
	self.FetcherGlobalStorage = dataStorage
	self.MetricStorage = dataStorage
	self.FetcherRunner = fetcherRunner
	self.DataControllerRunner = dataControllerRunner
	self.BlockchainSigner = pricingSigner
	//self.IntermediatorSigner = huoBiintermediatorSigner
	self.DepositSigner = depositSigner
	self.FeeBurnerAddress = burnerAddr
	self.NetworkAddress = networkAddr
	self.WhitelistAddress = whitelistAddr
	//self.ExchangeStorage = exsStorage
	// var huobiConfig common.HuobiConfig
	// exchangesIDs := os.Getenv("KYBER_EXCHANGES")
	// if strings.Contains(exchangesIDs, "huobi") {
	// 	huobiConfig = *self.GetHuobiConfig(kyberENV, addressConfig.Intermediator, huobiIntermediatorSigner)
	// }

	// create Exchange pool
	exchangePool := NewExchangePool(
		feeConfig,
		addressConfig,
		settingPath,
		self.Blockchain,
		minDeposit,
		kyberENV,
		self.Setting,
	)
	self.FetcherExchanges = exchangePool.FetcherExchanges()
	self.Exchanges = exchangePool.CoreExchanges()
}

var ConfigPaths = map[string]SettingPaths{
	mode.DEV_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://semi-node.kyber.network",
		[]string{
			"https://semi-node.kyber.network",
		},
		// "https://mainnet.infura.io",
		// []string{
		// 	"https://mainnet.infura.io",
		// },
	},
	mode.KOVAN_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://kovan.infura.io",
		[]string{},
	},
	mode.PRODUCTION_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_config.json",
		"https://mainnet.infura.io",
		[]string{
			"https://semi-node.kyber.network",
			"https://api.mycryptoapi.com/eth",
			"https://api.myetherapi.com/eth",
			"https://mew.giveth.io/",
		},
	},
	mode.MAINNET_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_config.json",
		"https://mainnet.infura.io",
		[]string{
			"https://mainnet.infura.io",
			"https://semi-node.kyber.network",
			"https://api.mycryptoapi.com/eth",
			"https://api.myetherapi.com/eth",
			"https://mew.giveth.io/",
		},
	},
	mode.STAGING_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_config.json",
		"https://mainnet.infura.io",
		[]string{
			"https://mainnet.infura.io",
			"https://semi-node.kyber.network",
			"https://api.mycryptoapi.com/eth",
			"https://api.myetherapi.com/eth",
			"https://mew.giveth.io/",
		},
	},
	mode.SIMULATION_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"http://blockchain:8545",
		[]string{
			"http://blockchain:8545",
		},
	},
	mode.ROPSTEN_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://ropsten.infura.io",
		[]string{
			"https://api.myetherapi.com/rop",
		},
	},
	mode.ANALYTIC_DEV_MODE: {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_analytics.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_fee_setrate.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"http://blockchain:8545",
		[]string{
			"http://blockchain:8545",
		},
	},
}

var Baseurl string = "http://127.0.0.1"

var BinanceInterfaces = make(map[string]binance.Interface)
var HuobiInterfaces = make(map[string]huobi.Interface)
var BittrexInterfaces = make(map[string]bittrex.Interface)

func SetInterface(base_url string) {
	BittrexInterfaces[mode.DEV_MODE] = bittrex.NewDevInterface()
	BittrexInterfaces[mode.KOVAN_MODE] = bittrex.NewKovanInterface(base_url)
	BittrexInterfaces[mode.MAINNET_MODE] = bittrex.NewRealInterface()
	BittrexInterfaces[mode.STAGING_MODE] = bittrex.NewRealInterface()
	BittrexInterfaces[mode.SIMULATION_MODE] = bittrex.NewSimulatedInterface(base_url)
	BittrexInterfaces[mode.ROPSTEN_MODE] = bittrex.NewRopstenInterface(base_url)
	BittrexInterfaces[mode.ANALYTIC_DEV_MODE] = bittrex.NewRopstenInterface(base_url)

	HuobiInterfaces[mode.DEV_MODE] = huobi.NewDevInterface()
	HuobiInterfaces[mode.KOVAN_MODE] = huobi.NewKovanInterface(base_url)
	HuobiInterfaces[mode.MAINNET_MODE] = huobi.NewRealInterface()
	HuobiInterfaces[mode.STAGING_MODE] = huobi.NewRealInterface()
	HuobiInterfaces[mode.SIMULATION_MODE] = huobi.NewSimulatedInterface(base_url)
	HuobiInterfaces[mode.ROPSTEN_MODE] = huobi.NewRopstenInterface(base_url)
	HuobiInterfaces[mode.ANALYTIC_DEV_MODE] = huobi.NewRopstenInterface(base_url)

	BinanceInterfaces[mode.DEV_MODE] = binance.NewDevInterface()
	BinanceInterfaces[mode.KOVAN_MODE] = binance.NewKovanInterface(base_url)
	BinanceInterfaces[mode.MAINNET_MODE] = binance.NewRealInterface()
	BinanceInterfaces[mode.STAGING_MODE] = binance.NewRealInterface()
	BinanceInterfaces[mode.SIMULATION_MODE] = binance.NewSimulatedInterface(base_url)
	BinanceInterfaces[mode.ROPSTEN_MODE] = binance.NewRopstenInterface(base_url)
	BinanceInterfaces[mode.ANALYTIC_DEV_MODE] = binance.NewRopstenInterface(base_url)
}
