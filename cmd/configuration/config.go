package configuration

import (
	"log"
	"os"
	"time"

	"github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/fetcher/http_runner"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/exchange/bittrex"
	"github.com/KyberNetwork/reserve-data/exchange/huobi"
	"github.com/KyberNetwork/reserve-data/http"
	"github.com/KyberNetwork/reserve-data/metric"
	"github.com/KyberNetwork/reserve-data/signer"
	"github.com/KyberNetwork/reserve-data/stat"
	statstorage "github.com/KyberNetwork/reserve-data/stat/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type SettingPaths struct {
	settingPath     string
	feePath         string
	dataStoragePath string
	statStoragePath string
	logStoragePath  string
	rateStoragePath string
	userStoragePath string
	signerPath      string
	endPoint        string
	bkendpoints     []string
}

type Config struct {
	ActivityStorage core.ActivityStorage
	DataStorage     data.Storage
	StatStorage     stat.StatStorage
	UserStorage     stat.UserStorage
	LogStorage      stat.LogStorage
	RateStorage     stat.RateStorage
	FetcherStorage  fetcher.Storage
	MetricStorage   metric.MetricStorage

	FetcherRunner     fetcher.FetcherRunner
	StatFetcherRunner stat.FetcherRunner
	FetcherExchanges  []fetcher.Exchange
	Exchanges         []common.Exchange
	BlockchainSigner  blockchain.Signer
	DepositSigner     blockchain.Signer

	EnableAuthentication bool
	AuthEngine           http.Authentication

	EthereumEndpoint        string
	BackupEthereumEndpoints []string

	SupportedTokens []common.Token

	WrapperAddress     ethereum.Address
	PricingAddress     ethereum.Address
	ReserveAddress     ethereum.Address
	FeeBurnerAddress   ethereum.Address
	NetworkAddress     ethereum.Address
	WhitelistAddress   ethereum.Address
	ThirdPartyReserves []ethereum.Address

	ChainType string
}

// GetStatConfig: load config to run stat server only
func (self *Config) AddStatConfig(setPath SettingPaths, addressConfig common.AddressConfig) {
	thirdpartyReserves := []ethereum.Address{}
	for _, address := range addressConfig.ThirdPartyReserves {
		thirdpartyReserves = append(thirdpartyReserves, ethereum.HexToAddress(address))
	}

	statStorage, err := statstorage.NewBoltStatStorage(setPath.statStoragePath)
	if err != nil {
		panic(err)
	}

	logStorage, err := statstorage.NewBoltLogStorage(setPath.logStoragePath)
	if err != nil {
		panic(err)
	}

	rateStorage, err := statstorage.NewBoltRateStorage(setPath.rateStoragePath)
	if err != nil {
		panic(err)
	}

	userStorage, err := statstorage.NewBoltUserStorage(setPath.userStoragePath)
	if err != nil {
		panic(err)
	}

	//fetcherRunner := http_runner.NewHttpRunner(8001)
	var statFetcherRunner stat.FetcherRunner

	if os.Getenv("KYBER_ENV") == "simulation" {
		statFetcherRunner = http_runner.NewHttpRunner(8002)
	} else {
		statFetcherRunner = fetcher.NewTickerRunner(7*time.Second, 5*time.Second, 3*time.Second, 5*time.Second, 5*time.Second, 10*time.Second, 7*time.Second)
	}

	self.StatStorage = statStorage
	self.UserStorage = userStorage
	self.LogStorage = logStorage
	self.RateStorage = rateStorage
	self.StatFetcherRunner = statFetcherRunner
	self.ThirdPartyReserves = thirdpartyReserves
}

func (self *Config) AddCoreConfig(setPath SettingPaths, authEnbl bool, addressConfig common.AddressConfig, kyberENV string) {
	networkAddr := ethereum.HexToAddress(addressConfig.Network)
	burnerAddr := ethereum.HexToAddress(addressConfig.FeeBurner)
	whitelistAddr := ethereum.HexToAddress(addressConfig.Whitelist)

	feeConfig, err := common.GetFeeFromFile(setPath.feePath)
	if err != nil {
		log.Fatalf("Fees file %s cannot found at: %s", setPath.feePath, err)
	}

	dataStorage, err := storage.NewBoltStorage(setPath.dataStoragePath)
	if err != nil {
		panic(err)
	}

	//fetcherRunner := http_runner.NewHttpRunner(8001)
	var fetcherRunner fetcher.FetcherRunner

	if os.Getenv("KYBER_ENV") == "simulation" {
		fetcherRunner = http_runner.NewHttpRunner(8001)
	} else {
		fetcherRunner = fetcher.NewTickerRunner(7*time.Second, 5*time.Second, 3*time.Second, 5*time.Second, 5*time.Second, 10*time.Second, 7*time.Second)
	}

	fileSigner, depositSigner := signer.NewFileSigner(setPath.signerPath)

	exchangePool := NewExchangePool(feeConfig, addressConfig, fileSigner, dataStorage, kyberENV)
	//exchangePool := exchangePoolFunc(feeConfig, addressConfig, fileSigner, storage)

	var hmac512auth http.KNAuthentication

	hmac512auth = http.KNAuthentication{
		fileSigner.KNSecret,
		fileSigner.KNReadOnly,
		fileSigner.KNConfiguration,
		fileSigner.KNConfirmConf,
	}

	if !authEnbl {
		log.Printf("\nWARNING: No authentication mode\n")
	}

	self.ActivityStorage = dataStorage
	self.DataStorage = dataStorage
	self.FetcherStorage = dataStorage
	self.MetricStorage = dataStorage
	self.FetcherRunner = fetcherRunner
	self.FetcherExchanges = exchangePool.FetcherExchanges()
	self.Exchanges = exchangePool.CoreExchanges()
	self.BlockchainSigner = fileSigner
	self.EnableAuthentication = authEnbl
	self.DepositSigner = depositSigner
	self.AuthEngine = hmac512auth
	self.FeeBurnerAddress = burnerAddr
	self.NetworkAddress = networkAddr
	self.WhitelistAddress = whitelistAddr
}

func (self *Config) MapTokens() map[string]common.Token {
	result := map[string]common.Token{}
	for _, t := range self.SupportedTokens {
		result[t.ID] = t
	}
	return result
}

var ConfigPaths = map[string]SettingPaths{
	"dev": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://semi-node.kyber.network",
		[]string{
			"https://semi-node.kyber.network",
		},
	},
	"kovan": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://kovan.infura.io",
		[]string{},
	},
	"production": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_config.json",
		"https://mainnet.infura.io/0BRKxQ0SFvAxGL72cbXi",
		[]string{
			"https://semi-node.kyber.network",
			"https://mainnet.infura.io",
			"https://api.mycryptoapi.com/eth",
			"https://api.myetherapi.com/eth",
			"https://mew.giveth.io/",
		},
	},
	"mainnet": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_config.json",
		"https://mainnet.infura.io/0BRKxQ0SFvAxGL72cbXi",
		[]string{
			"https://semi-node.kyber.network",
			"https://mainnet.infura.io",
			"https://api.mycryptoapi.com/eth",
			"https://api.myetherapi.com/eth",
			"https://mew.giveth.io/",
		},
	},
	"staging": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_config.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_users.db",
		"https://semi-node.kyber.network",
		[]string{
			"https://semi-node.kyber.network",
			"https://mainnet.infura.io",
			"https://api.mycryptoapi.com/eth",
			"https://api.myetherapi.com/eth",
			"https://mew.giveth.io/",
		},
	},
	"simulation": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/core_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"http://blockchain:8545",
		[]string{
			"http://blockchain:8545",
		},
	},
	"ropsten": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_stats.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_logs.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_rates.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/ropsten_users.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://ropsten.infura.io",
		[]string{
			"https://api.myetherapi.com/rop",
		},
	},
}

var Baseurl string = "http://127.0.0.1"

var BinanceInterfaces = make(map[string]binance.Interface)
var HuobiInterfaces = make(map[string]huobi.Interface)
var BittrexInterfaces = make(map[string]bittrex.Interface)

func SetInterface(base_url string) {

	BittrexInterfaces["dev"] = bittrex.NewDevInterface()
	BittrexInterfaces["kovan"] = bittrex.NewKovanInterface(base_url)
	BittrexInterfaces["mainnet"] = bittrex.NewRealInterface()
	BittrexInterfaces["staging"] = bittrex.NewRealInterface()
	BittrexInterfaces["simulation"] = bittrex.NewSimulatedInterface(base_url)
	BittrexInterfaces["ropsten"] = bittrex.NewRopstenInterface(base_url)

	HuobiInterfaces["dev"] = huobi.NewDevInterface()
	HuobiInterfaces["kovan"] = huobi.NewKovanInterface(base_url)
	HuobiInterfaces["mainnet"] = huobi.NewRealInterface()
	HuobiInterfaces["staging"] = huobi.NewRealInterface()
	HuobiInterfaces["simulation"] = huobi.NewSimulatedInterface(base_url)
	HuobiInterfaces["ropsten"] = huobi.NewRopstenInterface(base_url)

	BinanceInterfaces["dev"] = binance.NewDevInterface()
	BinanceInterfaces["kovan"] = binance.NewKovanInterface(base_url)
	BinanceInterfaces["mainnet"] = binance.NewRealInterface()
	BinanceInterfaces["staging"] = binance.NewRealInterface()
	BinanceInterfaces["simulation"] = binance.NewSimulatedInterface(base_url)
	BinanceInterfaces["ropsten"] = binance.NewRopstenInterface(base_url)
}

var HuobiAsync = map[string]bool{
	"dev":        false,
	"kovan":      true,
	"mainnet":    true,
	"staging":    true,
	"simulation": false,
	"ropsten":    true,
}
