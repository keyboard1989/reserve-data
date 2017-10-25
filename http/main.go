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
		ethereum.HexToAddress("0x7811f3b0505f621bac23cc0ad01bc8ccb68bbfdb"),
	)
	fetcher.AddExchange(exchange.NewLiqui())
	fetcher.AddExchange(exchange.NewBinance())
	fetcher.AddExchange(exchange.NewBittrex())
	fetcher.AddExchange(exchange.NewBitfinex())

	bc, err := blockchain.NewBlockchain(
		ethereum.HexToAddress("0x96aa24f61f16c28385e0a1c2ffa60a3518ded3ee"),
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
