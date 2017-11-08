package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/blockchain/nonce"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/core"
	corestorage "github.com/KyberNetwork/reserve-data/core/storage"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/exchange"
	"github.com/KyberNetwork/reserve-data/exchange/liqui"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	wrapperAddr := ethereum.HexToAddress("0x9624ec965ff9bf947a90d5b66392d41c53777c3d")
	reserveAddr := ethereum.HexToAddress("0xc9f8edc40f8b5369a3144bb29d7465b632fdb563")

	storage := storage.NewRamStorage()
	fetcherRunner := fetcher.NewTickerRunner(3*time.Second, 2*time.Second)
	fetcher := fetcher.NewFetcher(
		storage,
		fetcherRunner,
		reserveAddr,
	)

	fileSigner := signer.NewFileSigner("./config.json")

	// liqui := exchange.NewRealLiqui(fileSigner)
	liqui := exchange.NewLiqui(liqui.NewSimulatedLiquiEndpoint(fileSigner))
	common.SupportedExchanges[liqui.ID()] = liqui
	//fetcher.AddExchange(liqui)

	// binance
	binance := exchange.NewBinance(binance.NewSimulatedBinanceEndpoint(fileSigner))
	common.SupportedExchanges[binance.ID()] = binance
	fetcher.AddExchange(binance)

	// fetcher.AddExchange(exchange.NewBinance())
	// fetcher.AddExchange(exchange.NewBittrex())
	// fetcher.AddExchange(exchange.NewBitfinex())

	// endpoint := "http://localhost:8545"
	endpoint := "https://kovan.kyber.network"
	infura, err := ethclient.Dial(endpoint)
	if err != nil {
		panic(err)
	}

	// nonceCorpus := nonce.NewAutoIncreasing(infura, fileSigner)
	nonceCorpus := nonce.NewTimeWindow(infura, fileSigner)

	bc, err := blockchain.NewBlockchain(
		infura,
		wrapperAddr,
		reserveAddr,
		fileSigner,
		nonceCorpus,
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
		activityStorage := corestorage.NewRamStorage()
		core := core.NewReserveCore(bc, activityStorage, reserveAddr)
		server := NewHTTPServer(app, core, ":8000")
		server.Run()
	}
}
