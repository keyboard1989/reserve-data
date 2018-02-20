package fetcher

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Blockchain interface {
	FetchBalanceData(addr ethereum.Address, atBlock *big.Int, timepoint uint64) (map[string]common.BalanceEntry, error)
	// fetch current raw rates at specific block
	FetchRates(timepoint uint64, block uint64) (common.AllRateEntry, error)
	TxStatus(tx ethereum.Hash) (string, uint64, error)
	CurrentBlock() (uint64, error)
	SetRateMinedNonce() (uint64, error)
	GetLogs(fromBlock uint64, timepoint uint64) ([]common.TradeLog, error)
}
