package configuration

import (
	"log"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func GetAddressConfig(filePath string) common.AddressConfig {
	addressConfig, err := common.GetAddressConfigFromFile(filePath)
	if err != nil {
		log.Fatalf("Config file %s is not found. Check that KYBER_ENV is set correctly. Error: %s", filePath, err)
	}
	return addressConfig
}

func GetChainType(kyberENV string) string {
	switch kyberENV {
	case "mainnet", "production":
		return "byzantium"
	case "dev":
		return "homestead"
	case "kovan":
		return "homestead"
	case "staging":
		return "byzantium"
	case "simulation":
		return "homestead"
	case "ropsten":
		return "byzantium"
	default:
		return "homestead"
	}
}

func GetConfigPaths(kyberENV string) SettingPaths {
	switch kyberENV {
	case "mainnet", "production":
		return (ConfigPaths["mainnet"])
	case "dev":
		return (ConfigPaths["dev"])
	case "kovan":
		return (ConfigPaths["kovan"])
	case "staging":
		return (ConfigPaths["staging"])
	case "simulation":
		return (ConfigPaths["simulation"])
	case "ropsten":
		return (ConfigPaths["ropsten"])
	default:
		log.Println("Environment setting paths is not found, using dev...")
		return (ConfigPaths["dev"])
	}
}

func GetConfig(kyberENV string, authEnbl bool, endpointOW string, noCore, enableStat bool) *Config {
	setPath := GetConfigPaths(kyberENV)
	addressConfig := GetAddressConfig(setPath.settingPath)

	wrapperAddr := ethereum.HexToAddress(addressConfig.Wrapper)
	pricingAddr := ethereum.HexToAddress(addressConfig.Pricing)
	reserveAddr := ethereum.HexToAddress(addressConfig.Reserve)
	var intermediatorAddr ethereum.Address
	if addressConfig.Intermediator != "" {
		intermediatorAddr = ethereum.HexToAddress(addressConfig.Intermediator)
	} else {
		intermediatorAddr = ethereum.HexToAddress("0x13922F1857C0677F79e4BbB16Ad2c49fAa620829")
	}

	var endpoint string
	if endpointOW != "" {
		log.Printf("overwriting Endpoint with %s\n", endpointOW)
		endpoint = endpointOW
	} else {
		endpoint = setPath.endPoint
	}

	common.SupportedTokens = map[string]common.Token{}
	tokens := []common.Token{}
	for id, t := range addressConfig.Tokens {
		tok := common.Token{
			id, t.Address, t.Decimals,
		}
		common.SupportedTokens[id] = tok
		tokens = append(tokens, tok)
	}
	// log.Printf("Exchange storage is %s", setPath.exsStoragePath)
	// log.Printf("Stats    storage is %s", setPath.statStoragePath)
	// dataStorage, err := storage.NewBoltStorage(setPath.dataStoragePath)
	// if err != nil {
	// 	panic(err)
	// }
	// statStorage, err := statstorage.NewBoltStorage(setPath.statStoragePath)
	// if err != nil {
	// 	panic(err)
	// }
	// exsStorage, err := exsstorage.NewBoltStorage(setPath.exsStoragePath)
	// if err != nil {
	// 	panic(err)
	// }
	// //fetcherRunner := http_runner.NewHttpRunner(8001)
	// var fetcherRunner fetcher.FetcherRunner
	// var statFetcherRunner stat.FetcherRunner

	// if os.Getenv("KYBER_ENV") == "simulation" {
	// 	fetcherRunner = http_runner.NewHttpRunner(8001)
	// 	statFetcherRunner = http_runner.NewHttpRunner(8002)
	// } else {
	// 	fetcherRunner = fetcher.NewTickerRunner(10*time.Second, 4*time.Second, 11*time.Second, 13*time.Second, 15*time.Second)
	// 	statFetcherRunner = fetcher.NewTickerRunner(3*time.Second, 2*time.Second, 3*time.Second, 5*time.Second, 5*time.Second)
	// }
	// baseSigner := signer.GetBaseSigner(setPath.signerPath)
	// fileSigner := signer.NewFileSigner(baseSigner, baseSigner.Keystore, baseSigner.Passphrase)
	// depositSigner := signer.NewFileSigner(baseSigner, baseSigner.KeystoreD, baseSigner.PassphraseD)
	// intermediatorSigner := signer.NewFileSigner(baseSigner, baseSigner.KeystoreI, baseSigner.PassphraseI)

	// // endpoint := "https://ropsten.infura.io"
	// // endpoint := "http://blockchain:8545"
	// // endpoint := "https://kovan.infura.io"
	// var endpoint string
	// if endpointOW != "" {
	// 	log.Printf("overwriting Endpoint with %s\n", endpointOW)
	// 	endpoint = endpointOW
	// } else {
	// 	endpoint = setPath.endPoint
	// }

	// bkendpoints := setPath.bkendpoints
	// var hmac512auth http.KNAuthentication

	// hmac512auth = http.KNAuthentication{
	// 	fileSigner.KNSecret,
	// 	fileSigner.KNReadOnly,
	// 	fileSigner.KNConfiguration,
	// 	fileSigner.KNConfirmConf,
	// }
	// exchangePool := NewExchangePool(feeConfig, addressConfig, fileSigner, dataStorage, kyberENV, intermediatorSigner, endpoint, wrapperAddr, intermediatorAddr, exsStorage, authEnbl)
	// //exchangePool := exchangePoolFunc(feeConfig, addressConfig, fileSigner, storage)

	// if !authEnbl {
	// 	log.Printf("\nWARNING: No authentication mode\n")
	// }

	// chainType := GetChainType(kyberENV)

	// return &Config{
	// 	ActivityStorage:         dataStorage,
	// 	DataStorage:             dataStorage,
	// 	StatStorage:             statStorage,
	// 	FetcherStorage:          dataStorage,
	// 	ExchangeStorage:         exsStorage,
	// 	StatFetcherStorage:      statStorage,
	// 	MetricStorage:           dataStorage,
	// 	FetcherRunner:           fetcherRunner,
	// 	StatFetcherRunner:       statFetcherRunner,
	// 	FetcherExchanges:        exchangePool.FetcherExchanges(),
	// 	Exchanges:               exchangePool.CoreExchanges(),
	// 	BlockchainSigner:        fileSigner,
	// 	EnableAuthentication:    authEnbl,
	// 	DepositSigner:           depositSigner,
	// 	AuthEngine:              hmac512auth,

	bkendpoints := setPath.bkendpoints
	chainType := GetChainType(kyberENV)

	config := &Config{
		EthereumEndpoint:        endpoint,
		BackupEthereumEndpoints: bkendpoints,
		SupportedTokens:         tokens,
		WrapperAddress:          wrapperAddr,
		PricingAddress:          pricingAddr,
		ReserveAddress:          reserveAddr,
		ChainType:               chainType,
		IntermediatorAddress:    intermediatorAddr,
	}

	if enableStat {
		config.AddStatConfig(setPath, addressConfig)
	}

	if !noCore {
		config.AddCoreConfig(setPath, authEnbl, addressConfig, kyberENV)
	}
	return config
}
