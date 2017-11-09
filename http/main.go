package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func loadTimestamp(path string) []uint64 {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	timestamp := []uint64{}
	err = json.Unmarshal(raw, &timestamp)
	if err != nil {
		panic(err)
	}
	return timestamp
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	wrapperAddr := ethereum.HexToAddress("0x5aa7b0c53affef857523014ac6ce6c8d30bc68e6")
	reserveAddr := ethereum.HexToAddress("0x98990ee596d7c383a496f54c9e617ce7d2b3ed46")

	storage := storage.NewRamStorage()
	// fetcherRunner := fetcher.NewTickerRunner(3*time.Second, 2*time.Second)
	fetcherRunner := fetcher.NewTimestampRunner(
		loadTimestamp("/go/src/github.com/KyberNetwork/reserve-data/http/timestamps.json"),
		2*time.Second,
	)
	fetcher := fetcher.NewFetcher(
		storage,
		fetcherRunner,
		reserveAddr,
	)

	fileSigner := signer.NewFileSigner("/go/src/github.com/KyberNetwork/reserve-data/http/config.json")

	// liqui := exchange.NewRealLiqui(fileSigner)
	liqui := exchange.NewLiqui(liqui.NewSimulatedLiquiEndpoint(fileSigner))
	common.SupportedExchanges[liqui.ID()] = liqui
	fetcher.AddExchange(liqui)
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
