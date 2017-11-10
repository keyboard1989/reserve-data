package main

import (
	"github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Config struct {
	ActivityStorage core.ActivityStorage
	DataStorage     data.Storage
	FetcherStorage  fetcher.Storage

	FetcherRunner    fetcher.FetcherRunner
	FetcherExchanges []fetcher.Exchange
	Exchanges        []common.Exchange
	BlockchainSigner blockchain.Signer

	EthereumEndpoint string

	SupportedTokens []common.Token

	WrapperAddress ethereum.Address
	ReserveAddress ethereum.Address
}
