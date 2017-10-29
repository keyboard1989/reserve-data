package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/exchange"
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	storage := storage.NewRamStorage()
	fetcher := fetcher.NewFetcher(
		storage, 3*time.Second, 2*time.Second,
		ethereum.HexToAddress("0x7811f3b0505f621bac23cc0ad01bc8ccb68bbfdb"),
	)

	fileSigner := signer.NewFileSigner("config.json")

	liqui := exchange.NewLiqui(
		fileSigner,
		exchange.NewRealLiquiEndpoint(),
		// exchange.NewSimulatedLiquiEndpoint(),
	)
	common.SupportedExchanges[liqui.ID()] = liqui
	fetcher.AddExchange(liqui)
	// fetcher.AddExchange(exchange.NewBinance())
	// fetcher.AddExchange(exchange.NewBittrex())
	// fetcher.AddExchange(exchange.NewBitfinex())

	bc, err := blockchain.NewBlockchain(
		ethereum.HexToAddress("0x96aa24f61f16c28385e0a1c2ffa60a3518ded3ee"),
		fileSigner,
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
		app := data.NewReserveData(
			storage,
			fetcher,
		)
		app.Run()
		core := core.NewReserveCore(bc)
		server := NewHTTPServer(app, core, ":8000")
		server.Run()
	}
}
