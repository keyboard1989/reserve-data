package main

import (
	"github.com/KyberNetwork/reserve-data/alpha"
	"github.com/KyberNetwork/reserve-data/alpha/fetcher"
	"github.com/KyberNetwork/reserve-data/alpha/fetcher/blockchain"
	"github.com/KyberNetwork/reserve-data/alpha/fetcher/exchange"
	"github.com/KyberNetwork/reserve-data/alpha/storage"
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"

	"fmt"
	"runtime"
	"time"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	storage := storage.NewRamStorage()
	fetcher := fetcher.NewFetcher(
		storage, 3*time.Second, 2*time.Second,
		ethereum.HexToAddress("0x00f915055992d04e4ecde9c4724e4925142ea059"),
	)
	fetcher.AddExchange(exchange.NewLiqui())
	fetcher.AddExchange(exchange.NewBinance())
	fetcher.AddExchange(exchange.NewBittrex())
	fetcher.AddExchange(exchange.NewBitfinex())

	bc, err := blockchain.NewBlockchain(
		ethereum.HexToAddress("0x71c9Df0DDBa3Cc18383353e5B8A25B71c6fC397A"),
	)
	bc.AddToken(common.MustGetToken("ETH"))
	bc.AddToken(common.MustGetToken("OMG"))
	bc.AddToken(common.MustGetToken("DGD"))
	bc.AddToken(common.MustGetToken("CVC"))
	bc.AddToken(common.MustGetToken("MCO"))
	bc.AddToken(common.MustGetToken("GNT"))
	bc.AddToken(common.MustGetToken("ADX"))
	bc.AddToken(common.MustGetToken("EOS"))
	bc.AddToken(common.MustGetToken("PAY"))
	bc.AddToken(common.MustGetToken("BAT"))
	bc.AddToken(common.MustGetToken("KNC"))
	if err != nil {
		fmt.Printf("Can't connect to infura: %s\n", err)
	} else {
		fetcher.SetBlockchain(bc)
		app := alpha.NewReserveData(
			storage,
			fetcher,
		)
		app.Run()
		server := NewHTTPServer(app, ":8000")
		server.Run()
	}
}
