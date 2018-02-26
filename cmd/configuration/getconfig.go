package configuration

import (
	"log"
	"os"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/fetcher/http_runner"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/http"
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func GetAddressConfig(filePath string, addressOW [5]string) common.AddressConfig {
	addressConfig, err := common.GetAddressConfigFromFile(filePath)
	//addressConfig := GetAddressConfigFromViper()
	if err != nil {
		log.Fatalf("Config file %s is not found. Check that KYBER_ENV is set correctly. Error: %s", filePath, err)
	}
	if addressOW[0] != "" {
		log.Printf("Overwriting wrapper address config with %s \n", addressOW[0])
		addressConfig.Wrapper = addressOW[0]
	}
	if addressOW[1] != "" {
		log.Printf("Overwriting reserve address config with %s \n", addressOW[1])
		addressConfig.Wrapper = addressOW[0]
	}
	if addressOW[2] != "" {
		log.Printf("Overwriting pricing address config with %s \n", addressOW[2])
		addressConfig.Wrapper = addressOW[0]
	}
	if addressOW[3] != "" {
		log.Printf("Overwriting burner address config with %s \n", addressOW[3])
		addressConfig.Wrapper = addressOW[0]
	}
	if addressOW[4] != "" {
		log.Printf("Overwriting network address config with %s \n", addressOW[4])
		addressConfig.Wrapper = addressOW[0]
	}
	return addressConfig
}

// GetConfig: load and set all config with preset params and customize param depends on env
// This is to generalized all the getconfig function.
func GetConfig(kyberENV string, authEnbl bool, addressOW [5]string, endpointOW string) *Config {
	setPath := ConfigPaths[kyberENV]
	// settingPath := "/go/src/github.com/KyberNetwork/reserve-data/cmd/dev_setting.json"
	addressConfig := GetAddressConfig(setPath.settingPath, addressOW)

	feeConfig, err := common.GetFeeFromFile(setPath.feePath)
	if err != nil {
		log.Fatalf("Fees file %s cannot found at: %s", setPath.feePath, err)
	}
	wrapperAddr := ethereum.HexToAddress(addressConfig.Wrapper)
	reserveAddr := ethereum.HexToAddress(addressConfig.Reserve)
	pricingAddr := ethereum.HexToAddress(addressConfig.Pricing)
	burnerAddr := ethereum.HexToAddress(addressConfig.FeeBurner)
	networkAddr := ethereum.HexToAddress(addressConfig.Network)

	common.SupportedTokens = map[string]common.Token{}
	tokens := []common.Token{}
	for id, t := range addressConfig.Tokens {
		tok := common.Token{
			id, t.Address, t.Decimals,
		}
		common.SupportedTokens[id] = tok
		tokens = append(tokens, tok)
	}

	storage, err := storage.NewBoltStorage(setPath.storagePath)
	if err != nil {
		panic(err)
	}
	//fetcherRunner := http_runner.NewHttpRunner(8001)
	var fetcherRunner fetcher.FetcherRunner

	if os.Getenv("KYBER_ENV") == "simulation" {
		fetcherRunner = http_runner.NewHttpRunner(8001)
	} else {
		fetcherRunner = fetcher.NewTickerRunner(3*time.Second, 2*time.Second, 3*time.Second, 5*time.Second, 5*time.Second)
	}

	fileSigner, depositSigner := signer.NewFileSigner(setPath.signerPath)

	exchangePool := NewExchangePool(feeConfig, addressConfig, fileSigner, storage, kyberENV)
	//exchangePool := exchangePoolFunc(feeConfig, addressConfig, fileSigner, storage)

	// endpoint := "https://ropsten.infura.io"
	// endpoint := "http://blockchain:8545"
	// endpoint := "https://kovan.infura.io"
	var endpoint string
	if endpointOW != "" {
		log.Printf("overwriting Endpoint with %s\n", endpointOW)
		endpoint = endpointOW
	} else {
		endpoint = setPath.endPoint
	}

	bkendpoints := setPath.bkendpoints
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
	return &Config{
		ActivityStorage:         storage,
		DataStorage:             storage,
		FetcherStorage:          storage,
		MetricStorage:           storage,
		FetcherRunner:           fetcherRunner,
		FetcherExchanges:        exchangePool.FetcherExchanges(),
		Exchanges:               exchangePool.CoreExchanges(),
		BlockchainSigner:        fileSigner,
		EnableAuthentication:    authEnbl,
		DepositSigner:           depositSigner,
		AuthEngine:              hmac512auth,
		EthereumEndpoint:        endpoint,
		BackupEthereumEndpoints: bkendpoints,
		SupportedTokens:         tokens,
		WrapperAddress:          wrapperAddr,
		PricingAddress:          pricingAddr,
		ReserveAddress:          reserveAddr,
		FeeBurnerAddress:        burnerAddr,
		NetworkAddress:          networkAddr,
	}
}
