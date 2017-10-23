package main

import (
	"github.com/KyberNetwork/reserve-data/alpha"
	"github.com/KyberNetwork/reserve-data/alpha/fetcher"
	"github.com/KyberNetwork/reserve-data/alpha/fetcher/exchange"
	"github.com/KyberNetwork/reserve-data/alpha/storage"

	"time"
)

func main() {
	storage := storage.NewRamStorage()
	fetcher := fetcher.NewFetcher(storage, 3*time.Second)
	fetcher.AddExchange(exchange.NewLiqui())
	fetcher.AddExchange(exchange.NewBinance())
	fetcher.AddExchange(exchange.NewBittrex())
	fetcher.AddExchange(exchange.NewBitfinex())
	app := alpha.NewReserveData(
		storage,
		fetcher,
	)
	go app.Run()
	server := NewHTTPServer(app, ":8000")
	server.Run()
}
