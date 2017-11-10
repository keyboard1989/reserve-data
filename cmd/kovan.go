package main

import (
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	corestorage "github.com/KyberNetwork/reserve-data/core/storage"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/exchange"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/exchange/bitfinex"
	"github.com/KyberNetwork/reserve-data/exchange/bittrex"
	"github.com/KyberNetwork/reserve-data/exchange/liqui"
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func GetConfigForKovan() *Config {
	wrapperAddr := ethereum.HexToAddress("0x5aa7b0c53affef857523014ac6ce6c8d30bc68e6")
	reserveAddr := ethereum.HexToAddress("0x98990ee596d7c383a496f54c9e617ce7d2b3ed46")

	storage := storage.NewRamStorage()
	// storage, err := storage.NewBoltStorage("/go/src/github.com/KyberNetwork/reserve-data/cmd/core.db")
	// if err != nil {
	// 	panic(err)
	// }
	fetcherRunner := fetcher.NewTickerRunner(3*time.Second, 2*time.Second)
	// fetcherRunner := fetcher.NewTimestampRunner(
	// 	loadTimestamp("/go/src/github.com/KyberNetwork/reserve-data/cmd/timestamps.json"),
	// 	2*time.Second,
	// )

	fileSigner := signer.NewFileSigner("/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json")

	fetcherExchanges := []fetcher.Exchange{}
	// liqui := exchange.NewRealLiqui(fileSigner)
	liqui := exchange.NewLiqui(liqui.NewSimulatedLiquiEndpoint(fileSigner))
	bittrex := exchange.NewBittrex(bittrex.NewSimulatedBittrexEndpoint(fileSigner))
	binance := exchange.NewBinance(binance.NewSimulatedBinanceEndpoint(fileSigner))
	bitfinex := exchange.NewBitfinex(bitfinex.NewSimulatedBitfinexEndpoint(fileSigner))

	fetcherExchanges = append(fetcherExchanges, liqui)
	fetcherExchanges = append(fetcherExchanges, bittrex)
	fetcherExchanges = append(fetcherExchanges, binance)
	fetcherExchanges = append(fetcherExchanges, bitfinex)

	exchanges := []common.Exchange{}
	exchanges = append(exchanges, liqui)
	exchanges = append(exchanges, bittrex)
	exchanges = append(exchanges, binance)
	exchanges = append(exchanges, bitfinex)

	// endpoint := "http://localhost:8545"
	// endpoint := "https://kovan.kyber.network"
	endpoint := "https://kovan.infura.io"

	tokens := []common.Token{}
	tokens = append(tokens, common.MustGetToken("ETH"))
	tokens = append(tokens, common.MustGetToken("OMG"))
	tokens = append(tokens, common.MustGetToken("DGD"))
	tokens = append(tokens, common.MustGetToken("CVC"))
	tokens = append(tokens, common.MustGetToken("MCO"))
	tokens = append(tokens, common.MustGetToken("GNT"))
	tokens = append(tokens, common.MustGetToken("ADX"))
	tokens = append(tokens, common.MustGetToken("EOS"))
	tokens = append(tokens, common.MustGetToken("PAY"))
	tokens = append(tokens, common.MustGetToken("BAT"))
	tokens = append(tokens, common.MustGetToken("KNC"))

	activityStorage := corestorage.NewRamStorage()
	return &Config{
		ActivityStorage:  activityStorage,
		DataStorage:      storage,
		FetcherStorage:   storage,
		FetcherRunner:    fetcherRunner,
		FetcherExchanges: fetcherExchanges,
		Exchanges:        exchanges,
		BlockchainSigner: fileSigner,
		EthereumEndpoint: endpoint,
		SupportedTokens:  tokens,
		WrapperAddress:   wrapperAddr,
		ReserveAddress:   reserveAddr,
	}
}
