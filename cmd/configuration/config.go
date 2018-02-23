package configuration

import (
	"github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/exchange"
	"github.com/KyberNetwork/reserve-data/http"
	"github.com/KyberNetwork/reserve-data/metric"
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type SettingPaths struct {
	settingPath string
	feePath     string
	storagePath string
	signerPath  string
	endPoint    string
	bkendpoints []string
}

type Config struct {
	ActivityStorage core.ActivityStorage
	DataStorage     data.Storage
	FetcherStorage  fetcher.Storage
	MetricStorage   metric.MetricStorage

	FetcherRunner    fetcher.FetcherRunner
	FetcherExchanges []fetcher.Exchange
	Exchanges        []common.Exchange
	BlockchainSigner blockchain.Signer
	DepositSigner    blockchain.Signer

	EnableAuthentication bool
	AuthEngine           http.Authentication

	EthereumEndpoint        string
	BackupEthereumEndpoints []string

	SupportedTokens []common.Token

	WrapperAddress   ethereum.Address
	PricingAddress   ethereum.Address
	ReserveAddress   ethereum.Address
	FeeBurnerAddress ethereum.Address
	NetworkAddress   ethereum.Address
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
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/dev.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://mainnet.infura.io",
		[]string{
			"https://mainnet.infura.io",
		},
	},
	"kovan": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/kovan.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://kovan.infura.io",
		[]string{},
	},
	"mainnet": {
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/fee.json",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet.db",
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_config.json",
		"https://mainnet.infura.io",
		[]string{
			"https://node.kyber.network",
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
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/staging_config.json",
		"https://mainnet.infura.io",
		[]string{
			"https://node.kyber.network",
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
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json",
		"https://ropsten.infura.io",
		[]string{
			"https://api.myetherapi.com/rop",
		},
	},
}

var ExchangeFunction = map[string]func(common.ExchangeFeesConfig,
	common.AddressConfig,
	*signer.FileSigner,
	exchange.BittrexStorage) *ExchangePool{
	"dev":        NewDevExchangePool,
	"kovan":      NewKovanExchangePool,
	"mainnet":    NewMainnetExchangePool,
	"staging":    NewMainnetExchangePool,
	"simulation": NewSimulationExchangePool,
	"ropsten":    NewRopstenExchangePool,
}
