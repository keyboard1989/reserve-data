package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Blockchain interface {
	FetchBalanceData(addr ethereum.Address, timepoint uint64) (map[string]common.BalanceEntry, error)
	// fetch current raw rates at specific block
	FetchRates(timepoint uint64, block uint64) (common.AllRateEntry, error)
	TxStatus(tx ethereum.Hash) (string, error)
	CurrentBlock() (uint64, error)
	SetRateMinedNonce() (uint64, error)
	GetLogs(fromBlock uint64, timepoint uint64) ([]common.TradeLog, error)
	GetBlockNumByTxHash(ethereum.Hash) (uint64, error)
}
