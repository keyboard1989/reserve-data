package main

import (
	"log"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/exchange"
	// "github.com/KyberNetwork/reserve-data/exchange/binance"
	// "github.com/KyberNetwork/reserve-data/exchange/bitfinex"
	"github.com/KyberNetwork/reserve-data/exchange/bittrex"
	// "github.com/KyberNetwork/reserve-data/exchange/liqui"
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func GetConfigForMainnet() *Config {
	settingPath := "/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json"
	addressConfig, err := common.GetAddressConfigFromFile(settingPath)
	if err != nil {
		log.Fatalf("Config file %s is not found. Error: %s", settingPath, err)
	}
	wrapperAddr := ethereum.HexToAddress(addressConfig.Wrapper)
	reserveAddr := ethereum.HexToAddress(addressConfig.Reserve)

	common.SupportedTokens = map[string]common.Token{}
	tokens := []common.Token{}
	for id, t := range addressConfig.Tokens {
		tok := common.Token{
			id, t.Address, t.Decimals,
		}
		common.SupportedTokens[id] = tok
		tokens = append(tokens, tok)
	}

	// wrapperAddr := ethereum.HexToAddress("0x5aa7b0c53affef857523014ac6ce6c8d30bc68e6")
	// reserveAddr := ethereum.HexToAddress("0x98990ee596d7c383a496f54c9e617ce7d2b3ed46")

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
	// liqui := exchange.NewLiqui(liqui.NewKovanLiquiEndpoint(fileSigner))
	// for tokenID, addr := range addressConfig.Exchanges["liqui"] {
	// 	liqui.UpdateDepositAddress(common.MustGetToken(tokenID), addr)
	// }
	// binance := exchange.NewBinance(binance.NewRealBinanceEndpoint(fileSigner))
	// for tokenID, addr := range addressConfig.Exchanges["binance"] {
	// 	binance.UpdateDepositAddress(common.MustGetToken(tokenID), addr)
	// }
	bittrex := exchange.NewBittrex(bittrex.NewDevBittrexEndpoint(fileSigner))
	for tokenID, addr := range addressConfig.Exchanges["bittrex"] {
		bittrex.UpdateDepositAddress(common.MustGetToken(tokenID), addr)
	}
	// bitfinex := exchange.NewBitfinex(bitfinex.NewSimulatedBitfinexEndpoint(fileSigner))

	// fetcherExchanges = append(fetcherExchanges, liqui)
	// fetcherExchanges = append(fetcherExchanges, binance)
	fetcherExchanges = append(fetcherExchanges, bittrex)
	// fetcherExchanges = append(fetcherExchanges, bitfinex)

	exchanges := []common.Exchange{}
	// exchanges = append(exchanges, liqui)
	// exchanges = append(exchanges, binance)
	exchanges = append(exchanges, bittrex)
	// exchanges = append(exchanges, bitfinex)

	// endpoint := "http://localhost:8545"
	// endpoint := "https://kovan.kyber.network"
	endpoint := "https://mainnet.infura.io"

	return &Config{
		ActivityStorage:  storage,
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
