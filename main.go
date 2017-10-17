package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/apis"
	"github.com/KyberNetwork/reserve-data/fetcher"
	"github.com/KyberNetwork/reserve-data/fetcher/exchanges"
	"github.com/KyberNetwork/reserve-data/fetcher/tokens"
	"github.com/KyberNetwork/reserve-data/market"
	"github.com/KyberNetwork/reserve-data/storage"
)

func initLog() {
}

func main() {
	raven.SetDSN("https://bf15053001464a5195a81bc41b644751:eff41ac715114b20b940010208271b13@sentry.io/228067")
	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.GET("/prices", apis.AllPrices)
	r.GET("/prices/:base/:quote", apis.Price)

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Couldn't open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	market.StorageInstance = storage.NewRamStorage()
	market.GetPrice = market.GetCentralizedPrice
	market.GetAllPrice = market.GetAllCentralizedPrice

	fetcherInstance := fetcher.NewIntervalFetcher(
		market.StorageInstance,
		3*time.Second,
	)
	fetcherInstance.AddExchange(exchanges.Binance{})
	fetcherInstance.AddExchange(exchanges.Bittrex{})
	fetcherInstance.AddExchange(exchanges.Bitfinex{})
	fetcherInstance.AddExchange(exchanges.Liqui{})
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.OMG{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.ADX{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.BAT{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.CVC{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.DGD{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.EOS{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.FUN{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.GNT{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.KNC{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.LINK{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.MCO{},
			tokens.ETH{},
		},
	)
	fetcherInstance.AddTokenPair(
		fetcher.TokenPair{
			tokens.PAY{},
			tokens.ETH{},
		},
	)
	go fetcherInstance.Run()

	r.Run(":8000")
}
